# examples

to test each examples, first install the provider terraform-provider-osc
cd to the example you want to test

* edit terraform.tfvars with your informations : keypair, omi, current public IP ...

* plan and apply

'''
terraform plan && terraform apply
'''

* visualize

'''
terraform output
terraform show
'''

* destroy everything

'''
terraform destroy
'''
