resource "osc_customer_gateway" "customer_gw" {
  bgp_asn    = 65000
  ip_address = "149.6.164.226"
  type       = "ipsec.1"

  tags {
    project = "${var.project}"
    Name    = "${var.project}_customergw"
  }
}

