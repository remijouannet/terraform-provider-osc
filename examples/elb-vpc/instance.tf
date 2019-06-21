resource "osc_instance" "euw2-lb-1" {
  ami               = "${var.centos7}"
  availability_zone = "${var.region}a"
  instance_type     = "t2.micro"
  key_name          = "${var.sshkey}"
  subnet_id = "${osc_subnet.euw2-core-1-subnet-1.id}"

  vpc_security_group_ids = [
    "${osc_security_group.euw2-lb-1.id}",
  ]

  tags {
    Name    = "euw2-lb-1"
    project = "${var.project}"
  }
  depends_on = ["osc_elb.lb-1"]
}
