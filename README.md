Terraform Provider for Outscale
==================

Requirements
------------

-   [Terraform](https://www.terraform.io/downloads.html) 0.9.11 
-   [Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building it
---------------------

Download the binary and put it in the same folder than terraform binary

``sh
$ wget binary -o $(dirname $(which terraform))/terraform-provider-osc
$ 
```


Building it
---------------------

Clone repository to: `$GOPATH/src/github.com/remijouannet/terraform-provider-osc`

``sh
$ mkdir -p $GOPATH/src/github.com/remijouannet; cd $GOPATH/src/github.com/remijouannet
$ git clone git@github.com:remijouannet/terraform-provider-osc
``

Enter the provider directory and build the provider

``sh
$ cd $GOPATH/src/github.com/remijouannet/terraform-provider-osc
$ make build
```
