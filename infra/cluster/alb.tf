locals {
  alb_node_port = 30080
}

# public HTTP in, egress only to the nodes on the NodePort.
resource "aws_security_group" "alb" {
  name   = "${var.project}-alb"
  vpc_id = local.vpc_id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port       = local.alb_node_port
    to_port         = local.alb_node_port
    protocol        = "tcp"
    security_groups = [module.eks.node_security_group_id]
  }
}

# allow the ALB SG in on the NodePort
resource "aws_vpc_security_group_ingress_rule" "node_from_alb" {
  security_group_id            = module.eks.node_security_group_id
  from_port                    = local.alb_node_port
  to_port                      = local.alb_node_port
  ip_protocol                  = "tcp"
  referenced_security_group_id = aws_security_group.alb.id
}

resource "aws_lb" "app" {
  name               = "${var.project}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = local.subnet_ids
  idle_timeout       = 300   # comfortably above the 54s WS ping cadence
}

resource "aws_lb_target_group" "app" {
  name         = "${var.project}-tg"
  port         = local.alb_node_port
  protocol     = "HTTP"
  vpc_id       = local.vpc_id
  target_type  = "instance"

  health_check {
    path                = "/index.html"
    port                = local.alb_node_port
    protocol            = "HTTP"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    interval            = 10
    timeout             = 5
  }
}

resource "aws_lb_listener" "app" {
  load_balancer_arn = aws_lb.app.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app.arn
  }
}

# Register the node group's ASG so both nodes (and any future replacement)
# are auto-added to the target group.
resource "aws_autoscaling_attachment" "nodes" {
  autoscaling_group_name = module.eks.eks_managed_node_groups_autoscaling_group_names[0]
  lb_target_group_arn    = aws_lb_target_group.app.arn
}