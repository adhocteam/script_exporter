FROM golang:1.12.6-alpine AS build-env
MAINTAINER  James Kassemi (Ad Hoc, LLC) <james.kassemi@adhocteam.us>

RUN apk add --update git
RUN apk add --update gcc
RUN apk add --update libc-dev
RUN go get -u github.com/prometheus/promu

RUN mkdir script_exporter
COPY .promu.yml /go/script_exporter/
COPY script_exporter.go /go/script_exporter/
COPY go.mod /go/script_exporter/
COPY go.sum /go/script_exporter/

WORKDIR /go/script_exporter
RUN promu build

FROM alpine:3.8
COPY --from=build-env /go/script_exporter/script_exporter /bin/
COPY script-exporter.yml /etc/script-exporter/config.yml

EXPOSE      9172
ENTRYPOINT  [ "/bin/script-exporter" ]
CMD ["-config.file=/etc/script-exporter/config.yml"]
