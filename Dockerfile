FROM golang:1.13.1-alpine AS build-env

RUN apk add --update git gcc libc-dev
RUN go get -u github.com/prometheus/promu

RUN mkdir script_exporter
COPY .promu.yml script_exporter.go go.mod go.sum /go/script_exporter/

WORKDIR /go/script_exporter
RUN promu build

FROM alpine:3.8
LABEL upstream="https://github.com/adhocteam/script_exporter"
LABEL maintainer="james.kassemi@adhocteam.us"
RUN apk add --no-cache bash
COPY --from=build-env /go/script_exporter/script_exporter /bin/script-exporter
COPY script-exporter.yml /etc/script-exporter/config.yml

EXPOSE      9172
ENTRYPOINT  [ "/bin/script-exporter" ]
CMD ["-config.file=/etc/script-exporter/config.yml"]
