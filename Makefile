deps:
	go install github.com/hashicorp/terraform

build:
	go build -o terraform-provider-osc .

install:
	go install .

test:
	go test -v .
