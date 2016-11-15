package main

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
)

var config = &Config{
	Scripts: []*Script{
		{"success", "exit 0", 1},
		{"failure", "exit 1", 1},
		{"timeout", "sleep 5", 1},
	},
}

func TestRunScripts(t *testing.T) {
	runScripts(config)

	results := []struct {
		name    string
		total   int
		failure int
	}{
		{"success", 1, 0},
		{"failure", 1, 1},
		{"timeout", 1, 1},
	}

	for _, result := range results {
		s := &dto.Metric{}
		histogram.WithLabelValues(result.name).Write(s)

		if samples := int(*s.Histogram.SampleCount); samples != result.total {
			t.Errorf("Expecting 1 total sample, received %d", samples)
		}

		f := &dto.Metric{}
		failureHistogram.WithLabelValues(result.name).Write(f)

		if samples := int(*f.Histogram.SampleCount); samples != result.failure {
			t.Errorf("Expecting 1 failed sample, received %d", samples)
		}
	}
}
