resource "osc_security_group" "euw2-adm-1" {
  ingress {
    from_port   = 22
    to_port     = 22
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
    Name    = "euw2-adm-1"
    project = "${var.project}"
  }
}

resource "osc_security_group" "euw2-management-1" {
  ingress {
    from_port = 22
    to_port   = 22
    protocol  = "tcp"

    cidr_blocks = [
      "${osc_vpc.euw2-core-1-vpc-1.cidr_block}",
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  vpc_id = "${osc_vpc.euw2-core-1-vpc-1.id}"

  tags {
    Name    = "euw2-management-1"
    project = "${var.project}"
  }
}
