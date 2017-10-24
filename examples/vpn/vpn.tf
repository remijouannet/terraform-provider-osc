resource "osc_vpn_connection" "vpnconnection" {
  vpn_gateway_id      = "${osc_vpn_gateway.vpn_gw.id}"
  customer_gateway_id = "${osc_customer_gateway.customer_gw.id}"
  type                = "ipsec.1"
  static_routes_only  = false
}
