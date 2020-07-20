package main

import (
	"testing"
)

var config = &Config{
	Scripts: []*Script{
		{"success", "exit 0", 1},
		{"failure", "exit 1", 1},
		{"timeout", "sleep 5", 2},
	},
}

func TestRunScripts(t *testing.T) {
	measurements := runScripts(config.Scripts)

	expectedResults := map[string]struct {
		success     int
		minDuration float64
		exitCode	int
	}{
		"success": {1, 0, 0},
		"failure": {0, 0, 1},
		"timeout": {0, 2, -1},
	}

	for _, measurement := range measurements {
		expectedResult := expectedResults[measurement.Script.Name]

		if measurement.Success != expectedResult.success {
			t.Errorf("Expected result not found: %s", measurement.Script.Name)
		}

		if measurement.Duration < expectedResult.minDuration {
			t.Errorf("Expected duration %f < %f: %s", measurement.Duration, expectedResult.minDuration, measurement.Script.Name)
		}

		if measurement.ExitCode != expectedResult.exitCode {
			t.Errorf("Expected exit code %d != %d: %s", measurement.ExitCode, expectedResult.exitCode, measurement.Script.Name)
		}
	}
}

func TestScriptFilter(t *testing.T) {
	t.Run("RequiredParameters", func(t *testing.T) {
		_, err := scriptFilter(config.Scripts, "", "")

		if err.Error() != "`name` or `pattern` required" {
			t.Errorf("Expected failure when supplying no parameters")
		}
	})

	t.Run("NameMatch", func(t *testing.T) {
		scripts, err := scriptFilter(config.Scripts, "success", "")

		if err != nil {
			t.Errorf("Unexpected: %s", err.Error())
		}

		if len(scripts) != 1 || scripts[0] != config.Scripts[0] {
			t.Errorf("Expected script not found")
		}
	})

	t.Run("PatternMatch", func(t *testing.T) {
		scripts, err := scriptFilter(config.Scripts, "", "fail.*")

		if err != nil {
			t.Errorf("Unexpected: %s", err.Error())
		}

		if len(scripts) != 1 || scripts[0] != config.Scripts[1] {
			t.Errorf("Expected script not found")
		}
	})

	t.Run("AllMatch", func(t *testing.T) {
		scripts, err := scriptFilter(config.Scripts, "success", ".*")

		if err != nil {
			t.Errorf("Unexpected: %s", err.Error())
		}

		if len(scripts) != 3 {
			t.Fatalf("Expected 3 scripts, received %d", len(scripts))
		}

		for i, script := range config.Scripts {
			if scripts[i] != script {
				t.Fatalf("Expected script not found")
			}
		}
	})
}
