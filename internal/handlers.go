package internal

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (m *MemStorage) HandleMetric(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// fmt.Println("this is a post request")

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		fmt.Printf(
			"received metric with name: %s , type: %s , value: %s\n",
			metricName, metricType, metricValue,
		)

		if metricName == "" || metricValue == "" || metricType == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("metrics map now: %v\n", m)

		switch metricType {

		case "gauge":
			metricValueFloat, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				fmt.Printf("error parsing float: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			m.Gauge[metricName] = metricValueFloat
			w.WriteHeader(http.StatusOK)
			return

		case "counter":
			metricValueInt, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				fmt.Printf("error parsing int: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			m.Counter[metricName] += metricValueInt
			w.WriteHeader(http.StatusOK)
			return

		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
	w.WriteHeader(http.StatusBadRequest)
}
