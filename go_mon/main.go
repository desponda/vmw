package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	startAPIServer()
}

//external urls
var externalUrls = []string{
	"https://httpstat.us/200",
	"https://httpstat.us/503",
}

//staring server
func startAPIServer() {
	//registering the collectors to the prometheus
	for _, url := range externalUrls {
		req := newRequestCollector(url)
		prometheus.MustRegister(req)
	}

	logrus.Printf("setting up handler...")
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr: "localhost:8080",
	}
	logrus.Printf("listening on 8080")
	func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Fatalf("listening stopped: %v", err)
		}
	}()

}
