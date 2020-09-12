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

		assert.Equal(t, httpMetricsJSON["uri"], "/ping")
		assert.Equal(t, httpMetricsJSON["method"], "GET")
		assert.Equal(t, httpMetricsJSON["referer"], "https://www.example.com/ping")
		assert.Equal(t, httpMetricsJSON["userAgent"], "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
		assert.Equal(t, httpMetricsJSON["ip"], "217.73.23.163")
		assert.Equal(t, httpMetricsJSON["code"], 200)
		assert.Equal(t, httpMetricsJSON["duration"], 1)
	})
}
