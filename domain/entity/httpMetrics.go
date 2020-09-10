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
