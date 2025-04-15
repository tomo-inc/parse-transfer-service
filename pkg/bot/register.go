package bot

import (
	"context"
	"golang.org/x/time/rate"
)

var (
	notices = make([]Notificator, 0, 1)
)

func RegisterNotificator(n ...Notificator) {
	notices = append(notices, n...)
}

func SendMsg(ctx context.Context, msg Msg) {
	for _, notice := range notices {
		notice.Send(ctx, msg)
	}
}

func SendMsgWithRate(ctx context.Context, msg Msg, rate *rate.Limiter) {
	if rate.Allow() {
		for _, notice := range notices {
			notice.Send(ctx, msg)
		}
	}
}
