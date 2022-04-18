package main

import (
	"net/http"
	"testing"

	"bou.ke/monkey"
	"github.com/prometheus/client_golang/prometheus"
)

func TeststartAPIServer(t *testing.T) {
	monkey.Patch(newRequestCollector("xys"), func(url string) *RequestCollector {
		return &RequestCollector{
			ExternalUrl: prometheus.NewDesc("sample_external_url_up",
				"shows if url is up",
				nil, prometheus.Labels{
					"url": url},
			),
		}
	})
	monkey.Patch(http.Handle, func(pattern string, Handler http.Handler) {})
	defer monkey.UnpatchAll()
	startAPIServer()
}
