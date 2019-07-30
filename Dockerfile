FROM        quay.io/prometheus/busybox:latest
MAINTAINER  James Kassemi (Ad Hoc, LLC) <james.kassemi@adhocteam.us>

COPY script-exporter /bin/script-exporter
COPY script-exporter.yml /etc/script-exporter/config.yml
COPY runThis.sh /etc/script-exporter/runThis.sh

EXPOSE      9172
ENTRYPOINT  [ "/bin/script-exporter" ]
CMD ["-config.file=/etc/script-exporter/config.yml"]
