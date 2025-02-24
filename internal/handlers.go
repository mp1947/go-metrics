package internal

import (
	"fmt"
	"net/http"
	"strconv"
)

func (m *MemStorage) HandleMetric(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		metricType := r.PathValue("metricType")
		metricName := r.PathValue("metricName")
		metricValue := r.PathValue("metricValue")

		if metricName == "" || metricValue == "" || metricType == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("metric now: %v\n", m)

		switch metricType {

		case "gauge":
			metricValueFloat, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			m.Gauge[metricName] = metricValueFloat
			w.WriteHeader(http.StatusOK)
			return

		case "counter":
			metricValueInt, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			m.Counter[metricName] += metricValueInt
			w.WriteHeader(http.StatusOK)
			return

		default:
			w.WriteHeader(http.StatusBadRequest)
		}

	}

	w.WriteHeader(http.StatusBadRequest)

}
