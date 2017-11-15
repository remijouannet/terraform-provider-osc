resource "osc_security_group" "euw2-lb-1" {
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["${var.access_bastion}"]
  }

  tags {
    Name    = "euw2-lb-1"
    project = "${var.project}"
  }
}
