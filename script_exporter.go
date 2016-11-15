package main

import (
	"context"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	configFile    = flag.String("config.file", "script-exporter.yml", "Script exporter configuration file.")
	listenAddress = flag.String("web.listen-address", ":9172", "The address to listen on for HTTP requests.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	shell         = flag.String("config.shell", "/bin/sh", "Shell to execute script")

	histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "script_duration_seconds",
		Help: "Duration for configured scripts with zero exit status",
	}, []string{"script"})

	failureHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "script_failure_duration_seconds",
		Help: "Duration for configured scripts with non-zero exit status",
	}, []string{"script"})
)

type Config struct {
	Scripts []*Script `yaml:"scripts"`
}

type Script struct {
	Name    string `yaml:"name"`
	Content string `yaml:"script"`
	Timeout int64  `yaml:"timeout"`
}

func runScript(script *Script) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(script.Timeout)*time.Second)
	defer cancel()

	bashCmd := exec.CommandContext(ctx, *shell)

	bashIn, err := bashCmd.StdinPipe()

	if err != nil {
		return err
	}

	if err = bashCmd.Start(); err != nil {
		return err
	}

	if _, err = bashIn.Write([]byte(script.Content)); err != nil {
		return err
	}

	bashIn.Close()

	return bashCmd.Wait()
}

func runScripts(config *Config) {
	ch := make(chan bool)

	for _, script := range config.Scripts {
		go func(script *Script) {
			start := time.Now()
			err := runScript(script)
			duration := time.Since(start).Seconds()

			if err == nil {
				log.Debugf("OK: %s (after %fs).", script.Name, duration)
			} else {
				log.Errorf("ERROR: %s: %s (failed after %fs).", script.Name, err, duration)
				failureHistogram.WithLabelValues(script.Name).Observe(duration)
			}

			histogram.WithLabelValues(script.Name).Observe(duration)

			ch <- true
		}(script)
	}

	for i := 0; i < len(config.Scripts); i++ {
		<-ch
	}
}

func init() {
	prometheus.MustRegister(version.NewCollector("script_exporter"))
	prometheus.MustRegister(histogram)
	prometheus.MustRegister(failureHistogram)
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("script_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting script_exporter", version.Info())

	yamlFile, err := ioutil.ReadFile(*configFile)

	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	config := Config{}

	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}

	log.Infof("Loaded %d script configurations", len(config.Scripts))

	for _, script := range config.Scripts {
		if script.Timeout == 0 {
			script.Timeout = 15
		}
	}

	promHandler := prometheus.Handler()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		runScripts(&config)
		promHandler.ServeHTTP(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Script Exporter</title></head>
			<body>
			<h1>Script Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Infoln("Listening on", *listenAddress)

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
