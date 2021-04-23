package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/alecthomas/kong"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/thedustin/crushftp-prometheus-exporter/collector"
	"github.com/thedustin/crushftp-prometheus-exporter/crushftp"
)

var (
	Name         = "crushftp_exporter"
	ReadableName = "CrushFTP Prometheus Exporter"
)

var cli struct {
	ListenAddress    string   `env:"LISTEN_ADDRESS" help:"Address to listen on" short:"l" default:":9100"`
	MetricsEndpoint  string   `env:"METRICS_ENDPOINT" help:"Path under which to expose metrics" default:"/metrics"`
	CrushftpUrl      *url.URL `env:"CRUSHFTP_URL" help:"Base URL to the CrushFTP http(s) server" short:"H" default:"http://localhost"`
	CrushftpUsername string   `env:"CRUSHFTP_USERNAME" help:"CrushFTP admin username" short:"u"`
	CrushftpPassword string   `env:"CRUSHFTP_PASSWORD" help:"CrushFTP admin password" short:"p"`
	CrushftpInsecure bool     `env:"CRUSHFTP_INSECURE" help:"Ignore server certificate if using https"`
	Debug            bool     `help:"Enables debug mode and increases logging"`
	Version          bool     `help:"Display the application version" short:"V"`
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)

	kong.Parse(&cli)

	if cli.Debug {
		logger.Printf("Cli parameters: %+v", cli)
	}

	if cli.Version {
		fmt.Println(version.Print(ReadableName))
		return
	}

	clientOpts := crushftp.ClientOptions{
		HostAndPort: cli.CrushftpUrl.Host,
		Http:        crushftp.HttpClientOptions{Insecure: cli.CrushftpInsecure},
		Password:    cli.CrushftpPassword,
		PathBase:    cli.CrushftpUrl.Path,
		Scheme:      cli.CrushftpUrl.Scheme,
		Username:    cli.CrushftpUsername,
	}

	collectorOpts := collector.CollectorOpts{
		ErrorLogger: logger,
	}

	if cli.Debug {
		clientOpts.Logger = logger
		collectorOpts.Logger = logger
	}

	prometheus.MustRegister(version.NewCollector(Name))

	c := crushftp.NewClient(clientOpts)
	prometheus.MustRegister(collector.NewCollector(c, collectorOpts))

	http.Handle(cli.MetricsEndpoint, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", cli.MetricsEndpoint)

		w.Write([]byte(`<html>
             <head><title>` + ReadableName + `</title></head>
             <body>
             <h1>` + ReadableName + `</h1>
             <p><a href='` + cli.MetricsEndpoint + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	logger.Printf(
		"Starting exporter version %s at %q to collect data from CrushFTP at %q",
		version.Version,
		cli.ListenAddress,
		cli.CrushftpUrl,
	)

	logger.Fatal(http.ListenAndServe(cli.ListenAddress, nil))
}
