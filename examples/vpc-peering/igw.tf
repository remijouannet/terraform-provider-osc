resource "osc_internet_gateway" "euw2-core-1-igw-1" {
  vpc_id = "${osc_vpc.euw2-core-1-vpc-1.id}"

  tags {
    project = "${var.project}"
  }
}

resource "osc_internet_gateway" "euw2-core-2-igw-1" {
  vpc_id = "${osc_vpc.euw2-core-2-vpc-1.id}"

  tags {
    project = "${var.project}"
  }
}
