package bot

import (
	"context"
	"testing"
)

func TestLarkBot(t *testing.T) {
	bot := NewLarkBot("71e2e02e-e657-49b9-8e33-9fcfb51a3ce4", "")

	bot.Send(context.Background(), Msg{Title: "hahaha", Level: Info, Data: []*KeyValue{
		{
			Key:   "aaaa",
			Value: "dadsadad",
		}, {
			Key:   "bbbb",
			Value: "dadsadad",
		},
		{
			Key:   "cccc",
			Value: "cc",
		},
	}})
}
