provider "osc" {
  profile = "${var.profile}"
  region  = "${var.region}"

  endpoints {
    s3 = "${var.url_s3}"
  }
}
