package chain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/tomo.inc/parse-transfer-service/app/model"
	"github.com/tomo.inc/parse-transfer-service/pkg/bot"
	"github.com/tomo.inc/parse-transfer-service/share"
	"time"
)

type Tron struct {
	endpoint   string
	chainIndex string
	client     *ethclient.Client
}

func NewTron(endpoint string, chainIndex string) *Tron {
	client, err := ethclient.Dial(fmt.Sprintf("%s/jsonrpc", endpoint))
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to Ethereum client")
		bot.SendMsgWithRate(context.Background(), bot.Msg{
			Title: "connect failed",
			Level: bot.Error,
			Data: []*bot.KeyValue{
				{
					Key:   "endpoint",
					Value: endpoint,
				},
			},
		}, share.AlertLimiter)
	}
	cli := &Tron{
		endpoint:   endpoint,
		chainIndex: chainIndex,
		client:     client,
	}
	go cli.reConnect()
	return cli
}

func (t *Tron) reConnect() {
	tc := time.NewTicker(time.Second)
	for range tc.C {
		if t.client == nil {
			client, err := ethclient.Dial(fmt.Sprintf("%s/jsonrpc", t.endpoint))
			if err != nil {
				log.Error().Err(err).Msg("Failed to connect to Ethereum client")
				bot.SendMsgWithRate(context.Background(), bot.Msg{
					Title: "connect failed",
					Level: bot.Error,
					Data: []*bot.KeyValue{
						{
							Key:   "endpoint",
							Value: t.endpoint,
						},
					},
				}, share.AlertLimiter)
			} else {
				t.client = client
			}
		}
	}
}

func (t *Tron) FetchTxTransfer(tx string) ([]*model.Transfer, error) {
	if t.client != nil {
	}
	return nil, nil
}
