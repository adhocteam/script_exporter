# Script Exporter

GitHub: https://github.com/adhocteam/script_exporter

Prometheus exporter written to execute and collect metrics on script exit status
and duration. Designed to allow the execution of probes where support for the
probe type wasn't easily configured with the Prometheus blackbox exporter.

Minimum supported Go Version: 1.7.3

## Sample Configuration

```yaml
scripts:
  - name: success
    script: sleep 5

  - name: failure
    script: sleep 2 && exit 1

  - name: timeout
    script: sleep 5
    timeout: 1
```

## Running

You can run via docker with:

```
docker run -d -p 9172:9172 --name script-exporter \
  -v `pwd`/config.yml:/etc/script-exporter/config.yml:ro \
  -config.file=/etc/script-exporter/config.yml
  -web.listen-address=":9172" \
  -web.telemetry-path="/metrics" \
  -config.shell="/bin/sh" \
  adhocteam/script-exporter:master
```

You'll need to customize the docker image or use the binary on the host system
to install tools such as curl for certain scenarios.

## Output

`$ curl http://localhost:9172/metrics`

```
...
# HELP script_duration_seconds Duration for configured scripts with zero exit status
# TYPE script_duration_seconds histogram
script_duration_seconds_bucket{script="failure",le="0.005"} 0
script_duration_seconds_bucket{script="failure",le="0.01"} 0
script_duration_seconds_bucket{script="failure",le="0.025"} 0
script_duration_seconds_bucket{script="failure",le="0.05"} 0
script_duration_seconds_bucket{script="failure",le="0.1"} 0
script_duration_seconds_bucket{script="failure",le="0.25"} 0
script_duration_seconds_bucket{script="failure",le="0.5"} 0
script_duration_seconds_bucket{script="failure",le="1"} 0
script_duration_seconds_bucket{script="failure",le="2.5"} 1
script_duration_seconds_bucket{script="failure",le="5"} 1
script_duration_seconds_bucket{script="failure",le="10"} 1
script_duration_seconds_bucket{script="failure",le="+Inf"} 1
script_duration_seconds_sum{script="failure"} 2.014404206
script_duration_seconds_count{script="failure"} 1
script_duration_seconds_bucket{script="success",le="0.005"} 0
script_duration_seconds_bucket{script="success",le="0.01"} 0
script_duration_seconds_bucket{script="success",le="0.025"} 0
script_duration_seconds_bucket{script="success",le="0.05"} 0
script_duration_seconds_bucket{script="success",le="0.1"} 0
script_duration_seconds_bucket{script="success",le="0.25"} 0
script_duration_seconds_bucket{script="success",le="0.5"} 0
script_duration_seconds_bucket{script="success",le="1"} 0
script_duration_seconds_bucket{script="success",le="2.5"} 0
script_duration_seconds_bucket{script="success",le="5"} 0
script_duration_seconds_bucket{script="success",le="10"} 1
script_duration_seconds_bucket{script="success",le="+Inf"} 1
script_duration_seconds_sum{script="success"} 5.015765115
script_duration_seconds_count{script="success"} 1
script_duration_seconds_bucket{script="timeout",le="0.005"} 0
script_duration_seconds_bucket{script="timeout",le="0.01"} 0
script_duration_seconds_bucket{script="timeout",le="0.025"} 0
script_duration_seconds_bucket{script="timeout",le="0.05"} 0
script_duration_seconds_bucket{script="timeout",le="0.1"} 0
script_duration_seconds_bucket{script="timeout",le="0.25"} 0
script_duration_seconds_bucket{script="timeout",le="0.5"} 0
script_duration_seconds_bucket{script="timeout",le="1"} 0
script_duration_seconds_bucket{script="timeout",le="2.5"} 1
script_duration_seconds_bucket{script="timeout",le="5"} 1
script_duration_seconds_bucket{script="timeout",le="10"} 1
script_duration_seconds_bucket{script="timeout",le="+Inf"} 1
script_duration_seconds_sum{script="timeout"} 1.001585331
script_duration_seconds_count{script="timeout"} 1
# HELP script_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which script_exporter was built.
# TYPE script_exporter_build_info gauge
script_exporter_build_info{branch="",goversion="go1.7.3",revision="",version=""} 1
# HELP script_failure_duration_seconds Duration for configured scripts with non-zero exit status
# TYPE script_failure_duration_seconds histogram
script_failure_duration_seconds_bucket{script="failure",le="0.005"} 0
script_failure_duration_seconds_bucket{script="failure",le="0.01"} 0
script_failure_duration_seconds_bucket{script="failure",le="0.025"} 0
script_failure_duration_seconds_bucket{script="failure",le="0.05"} 0
script_failure_duration_seconds_bucket{script="failure",le="0.1"} 0
script_failure_duration_seconds_bucket{script="failure",le="0.25"} 0
script_failure_duration_seconds_bucket{script="failure",le="0.5"} 0
script_failure_duration_seconds_bucket{script="failure",le="1"} 0
script_failure_duration_seconds_bucket{script="failure",le="2.5"} 1
script_failure_duration_seconds_bucket{script="failure",le="5"} 1
script_failure_duration_seconds_bucket{script="failure",le="10"} 1
script_failure_duration_seconds_bucket{script="failure",le="+Inf"} 1
script_failure_duration_seconds_sum{script="failure"} 2.014404206
script_failure_duration_seconds_count{script="failure"} 1
script_failure_duration_seconds_bucket{script="timeout",le="0.005"} 0
script_failure_duration_seconds_bucket{script="timeout",le="0.01"} 0
script_failure_duration_seconds_bucket{script="timeout",le="0.025"} 0
script_failure_duration_seconds_bucket{script="timeout",le="0.05"} 0
script_failure_duration_seconds_bucket{script="timeout",le="0.1"} 0
script_failure_duration_seconds_bucket{script="timeout",le="0.25"} 0
script_failure_duration_seconds_bucket{script="timeout",le="0.5"} 0
script_failure_duration_seconds_bucket{script="timeout",le="1"} 0
script_failure_duration_seconds_bucket{script="timeout",le="2.5"} 1
script_failure_duration_seconds_bucket{script="timeout",le="5"} 1
script_failure_duration_seconds_bucket{script="timeout",le="10"} 1
script_failure_duration_seconds_bucket{script="timeout",le="+Inf"} 1
script_failure_duration_seconds_sum{script="timeout"} 1.001585331
script_failure_duration_seconds_count{script="timeout"} 1
```

## Design

The script exporter adds `script_failure_duration_seconds` and
`script_duration_seconds` histograms configured with the default buckets
to the default prometheus handler metrics available at `/metrics`. When Prometheus
scrapes the target all scripts are executed with sh and the observations are
included in the output.

YMMV if you're attempting to execute a large number of scripts, and you'd be
better off creating an exporter that can handle your protocol without launching
shell processes for each scrape.
