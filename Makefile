deps:
	go install github.com/hashicorp/terraform

build:
	go build -o terraform-provider-osc .

test:
	go test -v .
