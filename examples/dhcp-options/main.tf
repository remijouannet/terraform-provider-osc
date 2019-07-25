resource "osc_vpc" "vpc-1" {
    cidr_block = "192.168.0.0/16"
    tags {
        Name    = "${var.project}"
    }

}

resource "osc_subnet" "vpc-subnet-1" {
    vpc_id            = "${osc_vpc.vpc-1.id}"
    cidr_block        = "192.168.0.0/16"
    availability_zone = "eu-west-2a"

    tags {
        Name    = "${var.project}"
    }

}

resource "osc_internet_gateway" "gw-1" {
    vpc_id = "${osc_vpc.vpc-1.id}"
    tags {
        Name = "internet-gateway"
    }

}

resource "osc_route_table" "rt-1" {
    vpc_id          = "${osc_vpc.vpc-1.id}"
    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = "${osc_internet_gateway.gw-1.id}"
    }

    tags {
        Name = "Default-1"
    }

}

resource "osc_route_table_association" "attach-rt" {
    subnet_id = "${osc_subnet.vpc-subnet-1.id}"
    route_table_id = "${osc_route_table.rt-1.id}"
}

resource "osc_eip" "eip-1" {
    vpc = true
}

resource "osc_vpc_dhcp_options" "doptions-1" {
    domain_name = "contoso.lan"
    domain_name_servers = [ "192.168.0.10", "8.8.8.8", "8.8.4.4" ]
    ntp_servers = [ "129.250.35.251", "198.60.22.240", "138.68.46.177", "107.155.72.167" ]

}

resource "osc_vpc_dhcp_options_association" "attach-doptions-1" {
    vpc_id = "${osc_vpc.vpc-1.id}"
    dhcp_options_id = "${osc_vpc_dhcp_options.doptions-1.id}"
}
