provider "osc" {
  profile                     = "${var.profile}"
  region                      = "${var.region}"
  skip_region_validation      = true

  endpoints {
    ec2 = "${var.url_ec2}"
    iam = "${var.url_iam}"
  }
}
