resource "osc_subnet" "euw2-core-1-subnet-1" {
  vpc_id            = "${osc_vpc.euw2-core-1-vpc-1.id}"
  cidr_block        = "192.168.141.0/24"
  availability_zone = "${var.region}a"

  tags {
    Name    = "euw2-core-1-subnet-1"
    project = "${var.project}"
  }
}

resource "osc_subnet" "euw2-core-1-subnet-2" {
  vpc_id            = "${osc_vpc.euw2-core-1-vpc-1.id}"
  cidr_block        = "192.168.142.0/24"
  availability_zone = "${var.region}a"

  tags {
    Name    = "euw2-core-1-subnet-2"
    project = "${var.project}"
  }
}

resource "osc_subnet" "euw2-core-2-subnet-1" {
  vpc_id            = "${osc_vpc.euw2-core-2-vpc-1.id}"
  cidr_block        = "192.168.145.0/24"
  availability_zone = "${var.region}b"

  tags {
    Name    = "euw2-core-2-subnet-1"
    project = "${var.project}"
  }
}

resource "osc_subnet" "euw2-core-2-subnet-2" {
  vpc_id            = "${osc_vpc.euw2-core-2-vpc-1.id}"
  cidr_block        = "192.168.146.0/24"
  availability_zone = "${var.region}b"

  tags {
    Name    = "euw2-core-2-subnet-2"
    project = "${var.project}"
  }
}
