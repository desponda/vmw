package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

//response struct for external urls
type response struct {
	Url           string
	ResponseTime  float64
	ExternalUrlUp float64
}

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric
type RequestCollector struct {
	ExternalUrl  *prometheus.Desc
	ResponseTime *prometheus.Desc
}

//creating a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newRequestCollector(url string) *RequestCollector {
	return &RequestCollector{
		ExternalUrl: prometheus.NewDesc("sample_external_url_up",
			"shows if url is up",
			nil, prometheus.Labels{
				"url": url},
		),
		ResponseTime: prometheus.NewDesc("sample_external_url_response_ms",
			"response time for url",
			nil, prometheus.Labels{
				"url": url},
		),
	}
}

//Each and every collector must implement the Describe function.
func (collector *RequestCollector) Describe(ch chan<- *prometheus.Desc) {

	//Updating this section with the each metric created for a given collector
	ch <- collector.ExternalUrl
	ch <- collector.ResponseTime
}

//Collect implements required collect function for all promehteus collectors
func (collector *RequestCollector) Collect(ch chan<- prometheus.Metric) {

	//Implementing logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	resp := GetMetrics(collector.ExternalUrl.String())

	//Writing latest value for each metric in the prometheus metric channel.
	ch <- prometheus.MustNewConstMetric(collector.ExternalUrl, prometheus.CounterValue, resp.ExternalUrlUp)
	ch <- prometheus.MustNewConstMetric(collector.ResponseTime, prometheus.CounterValue, resp.ResponseTime)

}

//getting the results for the external urls using this method
func GetMetrics(url string) response {
	externalUrls := map[string]string{
		"200": "https://httpstat.us/200",
		"503": "https://httpstat.us/503",
	}
	var resp response
	if strings.Contains(url, externalUrls["200"]) {
		resp = handleHTTPRequest(externalUrls["200"])
	} else {
		resp = handleHTTPRequest(externalUrls["503"])
	}

	return resp

}

//this method runs the http resquest and get the response time for the request
func handleHTTPRequest(url string) response {

	start := time.Now()
	resp, err := http.Get(url)
	diff := time.Now().Sub(start).Milliseconds()
	if err != nil {
		logrus.Fatal(err)
	}
	result := response{
		Url:          url,
		ResponseTime: float64(diff),
	}
	if resp.StatusCode == 200 {
		result.ExternalUrlUp = 1
	}
	return result
}
