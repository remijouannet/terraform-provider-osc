resource "osc_ebs_volume" "example" {
    availability_zone = "${var.region}a"
    size = 40
    tags {
        Name = "HelloWorld"
        project = "${var.project}"
    }
}
