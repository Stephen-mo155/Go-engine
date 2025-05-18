package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Histogram of response time for handler.",
            Buckets: prometheus.DefBuckets,
        },
        []string{"handler", "method", "status"},
    )

    requestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"handler", "method", "status"},
    )
)

func init() {
    prometheus.MustRegister(requestDuration)
    prometheus.MustRegister(requestCount)
}

func handler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "2 %s", r.URL.Path[1:])

    duration := time.Since(start).Seconds()
    labels := prometheus.Labels{
        "handler": r.URL.Path,
        "method":  r.Method,
        "status":  "200",
    }
    requestDuration.With(labels).Observe(duration)
    requestCount.With(labels).Inc()
}

func main() {
    http.HandleFunc("/", handler)
    http.Handle("/metrics", promhttp.Handler())

    fmt.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}

