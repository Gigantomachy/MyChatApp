data "terraform_remote_state" "network" {
  backend = "local"
  config = {
    path = var.network_state_path
  }
}

locals {
  vpc_id     = data.terraform_remote_state.network.outputs.vpc_id
  subnet_ids = data.terraform_remote_state.network.outputs.public_subnet_ids
  az_letters = [for az in data.terraform_remote_state.network.outputs.azs : regex("[a-z]$", az)]
}

data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }
}

# SSH key for Ansible. Private key written to infra/db/mychatapp-key.pem (gitignored).
resource "tls_private_key" "mychatapp" {
  algorithm = "ED25519"
}

resource "aws_key_pair" "mychatapp" {
  key_name   = "${var.project}-db"
  public_key = tls_private_key.mychatapp.public_key_openssh
}

resource "local_file" "private_key" {
  content         = tls_private_key.mychatapp.private_key_openssh
  filename        = "${path.module}/mychatapp-key.pem"
  file_permission = "0600"
}

# Cassandra ports are open to the whole VPC CIDR so EKS pods (whose CNI IPs live
# inside the VPC) can reach the DB without a cross-layer SG dependency on the
# EKS node group.
resource "aws_security_group" "cassandra" {
  name   = "${var.project}-cassandra"
  vpc_id = local.vpc_id

  ingress {
    description = "CQL native transport"
    from_port   = 9042
    to_port     = 9042
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
  ingress {
    description = "inter-node gossip"
    from_port   = 7000
    to_port     = 7000
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
  ingress {
    description = "JMX (nodetool)"
    from_port   = 7199
    to_port     = 7199
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
  ingress {
    description = "encrypted client traffic (future TLS)"
    from_port   = 9142
    to_port     = 9142
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
  ingress {
    description = "SSH / Ansible"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.ssh_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "cassandra" {
  count                       = 3
  ami                         = data.aws_ami.al2023.id
  instance_type               = var.instance_type
  subnet_id                   = local.subnet_ids[count.index]
  private_ip                  = var.node_ips[count.index]
  vpc_security_group_ids      = [aws_security_group.cassandra.id]
  associate_public_ip_address = true   # Ansible/SSH from your machine; private IP is the stable address
  key_name                    = aws_key_pair.mychatapp.key_name

  root_block_device {
    volume_size = var.root_disk_size
    volume_type = "gp3"
    encrypted   = true
  }

  # data volume, formatted + mounted at /var/lib/cassandra by Ansible.
  # delete_on_termination = true: `db.sh destroy` wipes data (expected); use `db.sh stop` to keep it.
  ebs_block_device {
    device_name           = "/dev/sdf"
    volume_size           = var.data_disk_size
    volume_type           = "gp3"
    encrypted             = true
    delete_on_termination = true
  }

  tags = {
    Name = "${var.project}-db-${local.az_letters[count.index]}"
    Role = "cassandra"
    Rack = "rack-${local.az_letters[count.index]}"
  }
}