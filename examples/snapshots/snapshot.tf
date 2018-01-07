resource "osc_ebs_snapshot" "example_snapshot" {
  volume_id = "${osc_ebs_volume.example.id}"

  tags {
    Name = "HelloWorld_snap"
    project = "${var.project}"
  }
}
