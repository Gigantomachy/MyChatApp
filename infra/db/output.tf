output "instance_ids" {
  value = aws_instance.cassandra[*].id
}

output "private_ips" {
  value = aws_instance.cassandra[*].private_ip
}

output "public_ips" {
  value = aws_instance.cassandra[*].public_ip
}