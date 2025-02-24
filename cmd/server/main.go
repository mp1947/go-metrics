package main

import (
	"net/http"

	"github.com/mp1947/go-metrics/internal"
)

func main() {
	mux := http.NewServeMux()

	m := internal.MemStorage{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}

	mux.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", m.HandleMetric)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
