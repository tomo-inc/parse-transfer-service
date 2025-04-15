package share

import (
	"golang.org/x/time/rate"
	"time"
)

var (
	AlertLimiter = rate.NewLimiter(rate.Every(time.Minute), 1) // default
)

func SetAlertLimiter(duration time.Duration) {
	AlertLimiter = rate.NewLimiter(rate.Every(duration), 1)
}
