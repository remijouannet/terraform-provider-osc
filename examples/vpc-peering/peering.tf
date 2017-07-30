data "osc_caller_identity" "current" {}

resource "osc_vpc_peering_connection" "euw2-core-1-peering-1" {
  peer_owner_id = "${data.osc_caller_identity.current.account_id}"
  peer_vpc_id   = "${osc_vpc.euw2-core-2-vpc-1.id}"
  vpc_id        = "${osc_vpc.euw2-core-1-vpc-1.id}"
  auto_accept   = true

  tags {
    project = "${var.project}"
    Name    = "euw2-core-1-peering-1"
  }
}
