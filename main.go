package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	hello = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hello",
		Help: "Hello world count.",
	})
)

func main() {
	go func() {
		for {
			hello.Inc()
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
