TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
XC_ARCH="386 amd64 arm"
XC_OS="linux darwin windows freebsd openbsd solaris"
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"
VERSION=$$(git describe --abbrev=0 --tags)
PWD=$$(pwd)

COMMIT=$$(git rev-parse HEAD)
GOOS=$$(go env GOOS)
GOARCH=$$(go env GOARCH)

default: build

build: fmt
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

pkg: fmt
	mkdir -p ./pkg
	rm -rf ./pkg/*
	echo "==> Building..."
	CGO_ENABLED=0 gox -os=$(XC_OS) -arch=$(XC_ARCH) \
				-osarch=$(XC_EXCLUDE_OSARCH) \
				-output ./pkg/terraform-provider-osc_{{.OS}}_{{.Arch}}_$(VERSION)/terraform-provider-osc_$(VERSION) .

bin: fmt
	mkdir -p ./bin
	echo "==> Building..."
	CGO_ENABLED=0 gox -os=$(GOOS) -arch=$(GOARCH) -output ./bin/terraform-provider-osc_$(VERSION) .

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

release:
	bash scripts/github-releases.sh

docker-bin: docker-image
	docker run  \
		-v $(PWD)/bin:/go/src/github.com/remijouannet/terraform-provider-osc/bin \
		terraform-provider-osc:$(VERSION) bin

docker-image:
	docker build -t terraform-provider-osc:$(VERSION) .

docker-build:
	docker run \
		-v $(PWD)/pkg:/go/src/github.com/remijouannet/terraform-provider-osc/pkg \
		terraform-provider-osc:$(VERSION) pkg

docker-release:
	docker run \
		-v $(PWD)/pkg:/go/src/github.com/remijouannet/terraform-provider-osc/pkg \
		-e "GITHUB_TOKEN=$(GITHUB_TOKEN)" \
		terraform-provider-osc:$(VERSION) release 

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile
