package entity

type HTTPRequestMetrics struct {
	URI             string
	HTTPMethod      string
	Referer         string
	UserAgent       string
	IP              string
}

type HTTPResponseMetrics struct {
	Code    int
	Duration int
}
