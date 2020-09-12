package entity

type HTTPRequestMetrics struct {
	URI       string
	Method    string
	Referer   string
	UserAgent string
	IP        string
}

type HTTPResponseMetrics struct {
	Code     int
	Duration int
}

type HTTPMetrics struct {
	RequestMetrics  *HTTPRequestMetrics
	ResponseMetrics *HTTPResponseMetrics
}

func (httpMetrics *HTTPMetrics) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"uri":       httpMetrics.RequestMetrics.URI,
		"method":    httpMetrics.RequestMetrics.Method,
		"referer":   httpMetrics.RequestMetrics.Referer,
		"userAgent": httpMetrics.RequestMetrics.UserAgent,
		"ip":        httpMetrics.RequestMetrics.IP,
		"code":      httpMetrics.ResponseMetrics.Code,
		"duration":  httpMetrics.ResponseMetrics.Duration,
	}
}
