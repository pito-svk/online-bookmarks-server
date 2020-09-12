package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		requestMetrics := &HTTPRequestMetrics{
			URI:       "/ping",
			Method:    "GET",
			Referer:   "https://www.example.com/ping",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
			IP:        "217.73.23.163",
		}

		responseMetrics := &HTTPResponseMetrics{
			Code:     200,
			Duration: 1,
		}

		httpMetrics := &HTTPMetrics{
			RequestMetrics:  requestMetrics,
			ResponseMetrics: responseMetrics,
		}

		httpMetricsJSON := httpMetrics.ToMap()

		assert.Equal(t, "/ping", httpMetricsJSON["uri"])
		assert.Equal(t, "GET", httpMetricsJSON["method"])
		assert.Equal(t, "https://www.example.com/ping", httpMetricsJSON["referer"])
		assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36", httpMetricsJSON["userAgent"])
		assert.Equal(t, "217.73.23.163", httpMetricsJSON["ip"])
		assert.Equal(t, 200, httpMetricsJSON["code"])
		assert.Equal(t, 1, httpMetricsJSON["duration"])
	})
}
