# vpc-peering

this tf create two VPC with two subnets each, one public subnet and one private subnet,
a peering connection is set between the two VPC with the correct route table and security group
the output should show you the EIP of each adm and the private IP of each private instance,

to quickly test that the peering connection is working you can do the following

```
ssh -J root@IP.ADM.VPC.1 root@IP.PRIVATE.VPC.2 #use adm-1 on vpc-1 to access an instance on vpc-2
ssh -J root@IP.ADM.VPC.2 root@IP.PRIVATE.VPC.1 #use adm-2 on vpc-2 to access an instance on vpc-1
```
