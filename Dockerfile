FROM golang:1.13.3-alpine

MAINTAINER Rémi Jouannet "remijouannet@gmail.com"

RUN apk update
RUN apk add vim bash make git zip
RUN go get -u github.com/mitchellh/gox
RUN go get -u github.com/aktau/github-release

WORKDIR /go/src/github.com/remijouannet/terraform-provider-osc
COPY ./ /go/src/github.com/remijouannet/terraform-provider-osc

ENTRYPOINT ["make"]
