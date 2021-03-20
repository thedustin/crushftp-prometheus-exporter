package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/alecthomas/kong"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/thedustin/crushftp-prometheus-exporter/crushftp"
)

var (
	Name         = "crushftp_exporter"
	ReadableName = "Crushftp Prometheus Exporter"
)

var cli struct {
	ListenAddress    string   `env:"LISTEN_ADDRESS" help:"Address to listen on for web interface and telemetry" short:"l" default:":9100"`
	MetricsEndpoint  string   `env:"METRICS_ENDPOINT" help:"Path under which to expose metrics" default:"/metrics"`
	CrushftpUrl      *url.URL `env:"CRUSHFTP_URL" arg help:"URL to the CrushFTP http(s) server" default:"http://localhost"`
	CrushftpUsername string   `env:"CRUSHFTP_USERNAME" help:"Username for CrushFTP" short:"u"`
	CrushftpPassword string   `env:"CRUSHFTP_PASSWORD" help:"Username for CrushFTP" short:"p"`
	CrushftpInsecure bool     `env:"CRUSHFTP_INSECURE" help:"Ignore server certificate if using https" negatable`
	Debug            bool     `help:"Enables debug mode and increases logging"`
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	kong.Parse(&cli)

	logger.Printf("%+v", cli)

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector(Name))

	handler := promhttp.HandlerFor(
		r,
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		},
	)

	opts := crushftp.ClientOptions{
		HostAndPort: cli.CrushftpUrl.Host,
		Http:        crushftp.HttpClientOptions{Insecure: cli.CrushftpInsecure},
		Password:    cli.CrushftpPassword,
		PathBase:    cli.CrushftpUrl.Path,
		Scheme:      cli.CrushftpUrl.Scheme,
		Username:    cli.CrushftpUsername,
	}

	if cli.Debug {
		opts.Logger = logger
	}

	// c := crushftp.NewClient(opts)

	logger.Printf(
		"Starting exporter version %s at %q to collect data from CrushFTP at %q",
		version.Version,
		cli.ListenAddress,
		cli.CrushftpUrl,
	)

	http.Handle(cli.MetricsEndpoint, handler)
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

	logger.Fatal(http.ListenAndServe(cli.ListenAddress, nil))
}
