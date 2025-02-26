package main

import (
	"github.com/mp1947/go-metrics/internal"
)

func main() {

	metrics := internal.Metrics{
		Metric: make(map[string]interface{}),
	}

	go metrics.PollMetrics()
	go metrics.SendToServer()
	select {}
}
