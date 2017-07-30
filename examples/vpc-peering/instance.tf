resource "osc_instance" "euw2-adm-1" {
  ami               = "${var.centos7}"
  availability_zone = "${var.region}a"
  instance_type     = "t2.micro"
  key_name          = "${var.sshkey}"

  vpc_security_group_ids = [
    "${osc_security_group.euw2-adm-1.id}",
  ]

  subnet_id = "${osc_subnet.euw2-core-1-subnet-1.id}"

  tags {
    Name    = "euw2-adm-1"
    project = "${var.project}"
  }
}


resource "osc_instance" "euw2-jenkins-1" {
  ami               = "${var.centos7}"
  availability_zone = "${var.region}a"
  instance_type     = "t2.large"
  key_name          = "${var.sshkey}"

  vpc_security_group_ids = ["${osc_security_group.euw2-management-1.id}"]

  subnet_id = "${osc_subnet.euw2-core-1-subnet-2.id}"

  tags {
    Name    = "euw2-jenkins-1"
    project = "${var.project}"
  }
}

resource "osc_instance" "euw2-adm-2" {
  ami               = "${var.centos7}"
  availability_zone = "${var.region}b"
  instance_type     = "t2.micro"
  key_name          = "${var.sshkey}"

  vpc_security_group_ids = [
    "${osc_security_group.euw2-adm-2.id}",
  ]

  subnet_id = "${osc_subnet.euw2-core-2-subnet-1.id}"

  tags {
    Name    = "euw2-adm-2"
    project = "${var.project}"
  }
}

resource "osc_instance" "euw2-jenkins-2" {
  ami               = "${var.centos7}"
  availability_zone = "${var.region}b"
  instance_type     = "t2.micro"
  key_name          = "${var.sshkey}"

  vpc_security_group_ids = ["${osc_security_group.euw2-management-2.id}"]

  subnet_id = "${osc_subnet.euw2-core-2-subnet-1.id}"

  tags {
    Name    = "euw2-jenkins-2"
    project = "${var.project}"
  }
}
