resource "osc_eip" "euw2-core-1-eip-1" {
  vpc = true
}

resource "osc_eip" "euw2-core-1-eip-2" {
  network_interface = "${osc_instance.euw2-adm-1.network_interface_id}"
  vpc               = true
}
