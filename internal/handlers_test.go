package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleMetric(t *testing.T) {

	type request struct {
		method      string
		metricType  string
		path        string
		metricName  string
		metricValue string
	}

	tests := []struct {
		name         string
		request      request
		expectedCode int
	}{
		{
			name: "test bad request (get)",
			request: request{
				method:      http.MethodGet,
				metricType:  "type",
				metricName:  "name",
				path:        "/update",
				metricValue: "random",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "test correct request",
			request: request{
				method:      http.MethodPost,
				metricType:  "gauge",
				metricName:  "Alloc",
				metricValue: "1.1",
				path:        "/update",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "test incorrect metric type",
			request: request{
				method:      http.MethodPost,
				metricType:  "gaugeeeee",
				metricName:  "Alloc",
				metricValue: "1.1",
				path:        "/update",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	m := MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}

	mux := CreateMux(m)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := httptest.NewServer(mux)

			defer srv.Close()

			req := resty.New().R()
			req.Method = test.request.method
			req.URL = fmt.Sprintf(
				"%s%s/%s/%s/%s",
				srv.URL,
				test.request.path,
				test.request.metricType,
				test.request.metricName,
				test.request.metricValue,
			)

			t.Logf("sending %s request to %s", req.Method, req.URL)

			resp, err := req.Send()

			require.NoError(t, err, "error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode())
		})
	}
}
