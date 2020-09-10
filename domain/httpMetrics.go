package domain

import "peterparada.com/online-bookmarks/domain/entity"

type HTTPMetricsUsecase interface {
	GetHTTPRequestMetrics() entity.HTTPRequestMetrics
	GetHTTPResponseMetrics() entity.HTTPResponseMetrics
}
