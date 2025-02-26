package internal

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type Metrics struct {
	Metric map[string]interface{}
}

type (
	counter int64
	gauge   float64
)
