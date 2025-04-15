package bot

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func Test_rate(t *testing.T) {
	limiter := rate.NewLimiter(rate.Every(time.Second), 1)
	require.True(t, limiter.Allow())
	require.False(t, limiter.Allow())
}
