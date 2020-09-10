package http

import (
	"testing"

	"peterparada.com/online-bookmarks/domain/mocks"
)

func TestLogHTTPRequests(t *testing.T) {
	mockUsecase := mocks.NewHTTPMetricsUsecase()

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()

	})
}
