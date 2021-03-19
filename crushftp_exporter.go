package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/thedustin/crushftp-prometheus-exporter/crushftp"
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

	c := crushftp.NewClient("https", "example.org", "user", "pass")

	go func() {
		for {
			resp, err := c.GetDashboardItems()

			if err != nil {
				log.Fatal(err)
			} else {
				log.Printf("%+v", resp)

				time.Sleep(time.Second * 3)
			}
		}
	}()

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(":9100", nil))
}
