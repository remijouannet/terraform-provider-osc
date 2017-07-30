output "euw2-adm-1" {
  value = "${osc_instance.euw2-adm-1.public_ip}"
}

output "euw2-adm-2" {
  value = "${osc_instance.euw2-adm-2.public_ip}"
}

output "euw2-jenkins-1" {
  value = "${osc_instance.euw2-jenkins-1.private_ip}"
}

output "euw2-jenkins-2" {
  value = "${osc_instance.euw2-jenkins-2.private_ip}"
}
