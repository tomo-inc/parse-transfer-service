package chain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	apperr "github.com/tomo.inc/parse-transfer-service/app/err"
	"github.com/tomo.inc/parse-transfer-service/app/model"
	"github.com/tomo.inc/parse-transfer-service/pkg/bot"
	"github.com/tomo.inc/parse-transfer-service/pkg/contract/erc20"
	"github.com/tomo.inc/parse-transfer-service/pkg/task"
	"github.com/tomo.inc/parse-transfer-service/share"
	"time"
)

var (
	TransferTopic = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
)

type Evm struct {
	supportDebug bool
	endpoint     string
	chainIndex   string
	client       *ethclient.Client
}

func NewEvm(supportDebug bool, endpoint string, chainIndex string) *Evm {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to Ethereum client")
		bot.SendMsgWithRate(context.Background(), bot.Msg{
			Title: "evm connect failed",
			Level: bot.Error,
			Data: []*bot.KeyValue{
				{
					Key:   "chain index",
					Value: chainIndex,
				},
				{
					Key:   "endpoint",
					Value: endpoint,
				},
			},
		}, share.AlertLimiter)
	}
	cli := &Evm{
		supportDebug: supportDebug,
		endpoint:     endpoint,
		chainIndex:   chainIndex,
		client:       client,
	}
	go cli.reConnect()
	return cli
}

func (evm *Evm) reConnect() {
	tc := time.NewTicker(time.Second)
	for range tc.C {
		if evm.client == nil {
			client, err := ethclient.Dial(evm.endpoint)
			if err != nil {
				log.Error().Err(err).Msg("Failed to connect to Ethereum client")
				bot.SendMsgWithRate(context.Background(), bot.Msg{
					Title: "connect failed",
					Level: bot.Error,
					Data: []*bot.KeyValue{
						{
							Key:   "endpoint",
							Value: evm.endpoint,
						},
					},
				}, share.AlertLimiter)
			} else {
				evm.client = client
			}
		}
	}
}

func (e *Evm) FetchTxTransfer(tx string) (data []*model.Transfer, err error) {
	if len(common.FromHex(tx)) != common.HashLength {
		return nil, apperr.InvalidTxErr
	}

	if e.client != nil {
		logTransferTask := task.NewAsyncTask(context.Background(), fmt.Sprintf("%s-%s-log", e.chainIndex, tx), func() ([]*model.Transfer, error) {
			return e.parseLog(common.HexToHash(tx))
		})

		if e.supportDebug {
			debugTransfer := task.NewAsyncTask(context.Background(), fmt.Sprintf("%s-%s-debug", e.chainIndex, tx), func() ([]*model.Transfer, error) {
				return e.parseDebug(common.HexToHash(tx))
			})
			debugRes, err := debugTransfer.Result()
			if err != nil {
				return nil, err
			}
			data = append(data, debugRes...)
		}

		res, err := logTransferTask.Result()
		if err != nil {
			return nil, err
		}
		data = append(data, res...)
	}
	return data, nil
}

func (e *Evm) parseLog(tx common.Hash) (data []*model.Transfer, err error) {
	transfers := make([]*model.Transfer, 0)
	re, err := e.client.TransactionReceipt(context.Background(), tx)
	if err != nil {
		bot.SendMsg(context.Background(), bot.Msg{
			Title: "get receipt failed",
			Level: bot.Error,
			Data: []*bot.KeyValue{
				{
					Key:   "chain index",
					Value: e.chainIndex,
				},
				{
					Key:   "tx hash",
					Value: tx.String(),
				},
				{
					Key:   "err",
					Value: err.Error(),
				},
			},
		})
		log.Error().Str("tx", tx.Hex()).Err(err).Msg("Failed to get transaction receipt")
		return nil, err
	}
	for _, ethLog := range re.Logs {
		if len(ethLog.Topics) > 0 && ethLog.Topics[0].Cmp(TransferTopic) == 0 {
			transfer, err := erc20.TransferErc20Filter.ParseTransfer(*ethLog)
			if err != nil {
				continue
			}
			transfers = append(transfers, &model.Transfer{
				ChainIndex:   e.chainIndex,
				From:         transfer.From.String(),
				To:           transfer.To.String(),
				TokenAddress: ethLog.Address.String(),
				Amount:       transfer.Value.String(),
			})
		}
	}
	return transfers, nil
}

func (e *Evm) parseDebug(tx common.Hash) (data []*model.Transfer, err error) {
	var result = CallFrame{}
	err = e.client.Client().Call(&result, "debug_traceTransaction", tx.Hex(), &Config{Tracer: "callTracer"})
	if err != nil {
		bot.SendMsg(context.Background(), bot.Msg{
			Title: "get debug failed",
			Level: bot.Error,
			Data: []*bot.KeyValue{
				{
					Key:   "chain index",
					Value: e.chainIndex,
				},
				{
					Key:   "tx hash",
					Value: tx.String(),
				},
				{
					Key:   "err",
					Value: err.Error(),
				},
			},
		})
		return nil, err
	}

	return e.iterator(&result), nil
}

func (e *Evm) iterator(ett *CallFrame) []*model.Transfer {
	res := make([]*model.Transfer, 0)
	if ett != nil {
		if ett.Error == "" && ett.Type == "CALL" && ett.Value != "" && ett.Value != "0x0" {
			log.Debug().Str("type", ett.Type).Str("from", ett.From).Str("to", ett.To).Str("value", ett.Value).Msg("show debug")

			amount := hexutil.MustDecodeBig(ett.Value)
			res = append(res, &model.Transfer{
				From:         ett.From,
				To:           ett.To,
				Amount:       amount.String(),
				TokenAddress: zeroEvmAddress,
				ChainIndex:   e.chainIndex,
			})
		}

		for _, call := range ett.Calls {
			res = append(res, e.iterator(call)...)
		}
	}

	return res
}
