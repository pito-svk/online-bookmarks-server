package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalcDurationInMs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		now := time.Now()
		nowAfter55Ms := now.Add(time.Millisecond * 55)

		duration := calcDurationInMs(nowAfter55Ms, now)

		assert.Equal(t, 55, duration)
	})
}
