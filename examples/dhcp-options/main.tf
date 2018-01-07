resource "osc_vpc" "vpc-1" {
    cidr_block = "192.168.0.0/16"
    tags {
        project = "${var.project}"
    }
}

resource "osc_subnet" "vpc-subnet-1" {
    vpc_id            = "${osc_vpc.vpc-1.id}"
    cidr_block        = "192.168.1.0/24"
    availability_zone = "${var.region}a"

    tags {
        project = "${var.project}"
    }
}

resource "osc_vpc_dhcp_options" "doptions-1" {
    domain_name = "test.lan"
    domain_name_servers = [ "8.8.8.8", "8.8.4.4" ]
    ntp_servers = [ "8.8.8.8", "8.8.4.4" ]
}

resource "osc_vpc_dhcp_options_association" "attach-doptions-1" {
    vpc_id = "${osc_vpc.vpc-1.id}"
    dhcp_options_id = "${osc_vpc_dhcp_options.doptions-1.id}"
}
