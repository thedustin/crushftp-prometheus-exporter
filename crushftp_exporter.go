package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

var (
	Name = "crushftp_exporter"
)

func main() {
	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector(Name))

	handler := promhttp.HandlerFor(
		r,
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		},
	)

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(":9100", nil))
}
