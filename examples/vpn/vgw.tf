resource "osc_vpn_gateway" "vpn_gw" {
  vpc_id = "${osc_vpc.vpc.id}"
  tags {
    project = "${var.project}"
    Name    = "${var.project}_vpngw"
  }
}
