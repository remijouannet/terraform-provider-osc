resource "osc_elb" "lb-1" {
  name               = "${var.project}-prod"

  listener {
    instance_port      = 80
    instance_protocol  = "HTTP"
    lb_port            = 80
    lb_protocol        = "HTTP"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:80/"
    interval            = 30
  }

  security_groups = ["${osc_security_group.euw2-lb-1.id}"]
  subnets         = ["${osc_subnet.euw2-core-1-subnet-2.id}"]


  tags {
    Name = "${var.project}"
  }
}
