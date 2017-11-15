resource "osc_security_group" "euw2-lb-1" {
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["${var.access_bastion}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  vpc_id = "${osc_vpc.euw2-core-1-vpc-1.id}"

  tags {
    Name    = "euw2-lb-1"
    project = "${var.project}"
  }
}
