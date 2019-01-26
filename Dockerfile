FROM golang:1.11.5-alpine

MAINTAINER Rémi Jouannet "remijouannet@gmail.com"

RUN apk update
RUN apk add bash make git zip
RUN go get -u github.com/kardianos/govendor
RUN go get -u github.com/mitchellh/gox
RUN go get -u github.com/aktau/github-release

WORKDIR /go/src/github.com/remijouannet/terraform-provider-osc
COPY ./ /go/src/github.com/remijouannet/terraform-provider-osc

ENTRYPOINT ["make"]
