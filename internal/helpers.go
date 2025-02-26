package internal

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	reportInterval = 10
	pollInterval   = 2
	gaugeMetric    = "gauge"
	counterMetric  = "counter"
)

func (m *Metrics) SendToServer() {
	for {
		fmt.Println("sending metrics...")
		for k, v := range m.Metric {
			serverURL := fmt.Sprintf(
				"http://localhost:8080/update/%s/%s/%v",
				reflect.TypeOf(v).Name(), k, v)

			resp, err := http.Post(serverURL, "text/plain", nil)

			if err != nil {
				panic(err)
			}

			switch resp.StatusCode {
			case http.StatusOK:
				resp.Body.Close()
				continue
			default:
				errorMessage := fmt.Sprintf("non 200 response code: %d", resp.StatusCode)
				resp.Body.Close()
				panic(errors.New(errorMessage))
			}
		}
		time.Sleep(time.Second * reportInterval)
	}

}

func (m Metrics) PollMetrics() {
	var ms runtime.MemStats
	pollCount := 0

	for {

		fmt.Println("polling metrics")
		runtime.ReadMemStats(&ms)
		m.Metric["PollCount"] = counter(pollCount)
		m.Metric["RandomValue"] = gauge(rand.Float64())

		m.Metric["Alloc"] = gauge(ms.Alloc)
		m.Metric["BuckHashSys"] = gauge(ms.BuckHashSys)
		m.Metric["GCCPUFraction"] = gauge(ms.GCCPUFraction)
		m.Metric["GCSys"] = gauge(ms.GCSys)
		m.Metric["HeapAlloc"] = gauge(ms.HeapAlloc)
		m.Metric["HeapIdle"] = gauge(ms.HeapIdle)
		m.Metric["HeapInuse"] = gauge(ms.HeapInuse)
		m.Metric["HeapObjects"] = gauge(ms.HeapObjects)
		m.Metric["HeapReleased"] = gauge(ms.HeapReleased)
		m.Metric["HeapSys"] = gauge(ms.HeapSys)
		m.Metric["LastGC"] = gauge(ms.LastGC)
		m.Metric["Lookups"] = gauge(ms.Lookups)
		m.Metric["MCacheInuse"] = gauge(ms.MCacheInuse)
		m.Metric["MCacheSys"] = gauge(ms.MCacheSys)
		m.Metric["MSpanInuse"] = gauge(ms.MSpanInuse)
		m.Metric["MSpanSys"] = gauge(ms.MSpanSys)
		m.Metric["Mallocs"] = gauge(ms.Mallocs)
		m.Metric["NextGC"] = gauge(ms.NextGC)
		m.Metric["NumForcedGC"] = gauge(ms.NumForcedGC)
		m.Metric["NumGC"] = gauge(ms.NumGC)
		m.Metric["OtherSys"] = gauge(ms.OtherSys)
		m.Metric["PauseTotalNs"] = gauge(ms.PauseTotalNs)
		m.Metric["StackInuse"] = gauge(ms.StackInuse)
		m.Metric["StackSys"] = gauge(ms.StackSys)
		m.Metric["Sys"] = gauge(ms.Sys)
		m.Metric["TotalAlloc"] = gauge(ms.TotalAlloc)

		pollCount += 1
		time.Sleep(time.Second * pollInterval)
	}
}

func CreateRouter(m MemStorage) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Handle("/{metricType}/{metricName}/{metricValue}", http.HandlerFunc(m.HandleUpdateMetric))
	})

	r.Route("/value", func(r chi.Router) {
		r.Handle("/{metricType}/{metricName}", http.HandlerFunc(m.HandleGetMetric))
	})
	return r
}
