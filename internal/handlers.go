package internal

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (m *MemStorage) HandleMetric(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	fmt.Printf(
		"received metric with name: %s , type: %s \n",
		metricName, metricType,
	)

	if metricName == "" || metricType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		metricValue := chi.URLParam(r, "metricValue")

		if metricValue == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch metricType {

		case gaugeMetric:
			metricValueFloat, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				fmt.Printf("error parsing float: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			m.Gauge[metricName] = metricValueFloat
			w.WriteHeader(http.StatusOK)
			return

		case counterMetric:
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

	case http.MethodGet:
		switch metricType {
		case gaugeMetric:
			v := m.Gauge[metricName]
			if v == float64(0) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%f", v)))
			return
		case counterMetric:
			v := m.Counter[metricName]
			if v > 0 {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%d", v)))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
