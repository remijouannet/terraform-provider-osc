resource "osc_vpc" "euw2-core-1-vpc-1" {
  cidr_block = "192.168.140.0/22"

  tags {
    project = "${var.project}"
    Name    = "euw2-core-1-vpc-1"
  }
}
