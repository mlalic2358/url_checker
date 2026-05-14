package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var(
	urlUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "url_up",
			Help: "Whether the URL is up(1) or down(0)",
		},
		[]string{"url"},
	)
	urlLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "url_check_latency_seconds",
			Help: "Latency of URLs in seconds",
		},
		[]string{"url"},
	)
)

type CheckResult struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Error      error
}

func checkURL(url string) CheckResult {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return CheckResult{URL: url, Error: err}
	}
	defer resp.Body.Close()

	return CheckResult{URL: url, StatusCode: resp.StatusCode, Latency: time.Since(start)}
}

func main() {

	urls := []string{
		"https://github.com",
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/500",
		"https://httpbin.org/delay/2",
		"https://google.com",
	}
	
    prometheus.MustRegister(urlUp)
    prometheus.MustRegister(urlLatency)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	ticker := time.NewTicker(30 * time.Second)

	for range ticker.C {
			fmt.Printf("%s\n", time.Now())

		for _, url := range urls {
			result := checkURL(url)
			if result.Error != nil {
				urlUp.WithLabelValues(url).Set(0)
				// fmt.Printf("%-40s ERROR: %v\n", result.URL, result.Error)
			} else {
					if result.StatusCode >= 500 {
						urlUp.WithLabelValues(url).Set(0)
					} else {
						urlUp.WithLabelValues(url).Set(1)
					}
					urlLatency.WithLabelValues(url).Observe(result.Latency.Seconds())
				// fmt.Printf("%-40s %d %s\n", result.URL, result.StatusCode, result.Latency)
			}
		}
	}

}
