resource "osc_route_table" "euw2-core-1-route-1" {
  vpc_id = "${osc_vpc.euw2-core-1-vpc-1.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${osc_internet_gateway.euw2-core-1-igw-1.id}"
  }

  tags {
    Name    = "euw2-core-1-route-1"
    project = "${var.project}"
  }
}

resource "osc_route_table" "euw2-core-1-route-2" {
  vpc_id = "${osc_vpc.euw2-core-1-vpc-1.id}"

  route {
    cidr_block = "0.0.0.0/0"
    nat_gateway_id = "${osc_nat_gateway.euw2-core-1-natgw-1.id}"
  }

  tags {
    Name    = "euw2-core-1-route-2"
    project = "${var.project}"
  }
}

resource "osc_route_table_association" "euw2-core-1-route-1" {
  subnet_id      = "${osc_subnet.euw2-core-1-subnet-1.id}"
  route_table_id = "${osc_route_table.euw2-core-1-route-1.id}"
}

resource "osc_route_table_association" "euw2-core-1-route-2" {
  subnet_id      = "${osc_subnet.euw2-core-1-subnet-2.id}"
  route_table_id = "${osc_route_table.euw2-core-1-route-2.id}"
}
