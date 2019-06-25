resource "osc_instance" "euw2-lb-1" {
  ami               = "${var.centos7}"
  availability_zone = "${var.region}a"
  instance_type     = "t2.micro"
  key_name          = "${var.sshkey}"
   depends_on = ["osc_elb.lb-1"]
  security_groups = ["${osc_security_group.euw2-lb-1.id}"]

  tags {
    Name    = "euw2-lb-1"
    project = "${var.project}"
  }
}
