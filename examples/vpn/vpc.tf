resource "osc_vpc" "vpc" {
  cidr_block = "192.168.140.0/22"

  tags {
    project = "${var.project}"
    Name    = "${var.project}_vpc"
  }
}
