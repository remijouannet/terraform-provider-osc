resource "osc_elb_attachment" "attach-1" {
   elb = "${osc_elb.lb-1.id}"
   instance = "${osc_instance.euw2-lb-1.id}"
}
