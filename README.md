[![Github All Releases](https://img.shields.io/github/downloads/remijouannet/terraform-provider-osc/total.svg)]()

Terraform Provider for Outscale (unofficial)
==================

Requirements
------------

-   [Terraform](https://www.terraform.io/downloads.html) 0.11.11 
-   [Go](https://golang.org/doc/install) 1.11.5 (to build the provider plugin)

Install
---------------------

Download the binary [here](https://github.com/remijouannet/terraform-provider-osc/releases)

the documentation for [installing Third-party plugins](https://www.terraform.io/docs/plugins/basics.html#installing-plugins)

Build without docker
---------------------

Clone repository to: `$GOPATH/src/github.com/remijouannet/terraform-provider-osc`

```
$ mkdir -p $GOPATH/src/github.com/remijouannet; cd $GOPATH/src/github.com/remijouannet
$ git clone git@github.com:remijouannet/terraform-provider-osc
```

Enter the provider directory and build the provider

```
$ cd $GOPATH/src/github.com/remijouannet/terraform-provider-osc
$ make build
```

Build with docker
---------------------

build the docker image

```
$ make docker-image
```

build the binaries, you'll find all the binaries in pkg/

```
$ make docker-build
```
