package chain

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/mr-tron/base58"
	"github.com/rs/zerolog/log"
	"github.com/tomo.inc/parse-transfer-service/app/model"
	"github.com/tomo.inc/parse-transfer-service/pkg/bot"
	"github.com/tomo.inc/parse-transfer-service/share"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"math/big"
	"strconv"
	"strings"
)

type Tron struct {
	endpoint, token string
	chainIndex      string
	client          *client.GrpcClient
}

func NewTron(endpoint, token, chainIndex string) *Tron {
	client := client.NewGrpcClient(endpoint)
	err := client.Start(grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithPerRPCCredentials(&auth{
			token: token,
		}))
	if err != nil {
		bot.SendMsgWithRate(context.Background(), bot.Msg{
			Title: "tron connect failed",
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
		log.Fatal().Err(err).Msg("Failed to connect to Tron client")
	}
	return &Tron{
		endpoint:   endpoint,
		token:      token,
		chainIndex: chainIndex,
		client:     client,
	}
}

func (t *Tron) FetchTxTransfer(tx string) ([]*model.Transfer, error) {
	tx = strings.TrimPrefix(tx, "0x")
	txRes, err := t.client.GetTransactionInfoByID(tx)
	if err != nil {
		log.Error().Str("tx", tx).Err(err).Msg("Failed to get transaction receipt")
		bot.SendMsg(context.Background(), bot.Msg{
			Title: "get receipt failed",
			Level: bot.Error,
			Data: []*bot.KeyValue{
				{
					Key:   "chain index",
					Value: t.chainIndex,
				},
				{
					Key:   "tx hash",
					Value: tx,
				},
				{
					Key:   "err",
					Value: err.Error(),
				},
			},
		})
		return nil, err
	}
	transfers := make([]*model.Transfer, 0, 1)
	if len(txRes.Log) > 0 {
		for _, tronLog := range txRes.Log {
			if len(tronLog.Topics) == 3 && len(tronLog.Data) == 32 && bytes.Equal(tronLog.Topics[0], common.Hex2Bytes("ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")) {
				transfers = append(transfers, &model.Transfer{
					ChainIndex:   t.chainIndex,
					TokenAddress: encodeTronAddr(tronLog.Address),
					From:         encodeTronAddr(tronLog.Topics[1][12:]),
					To:           encodeTronAddr(tronLog.Topics[2][12:]),
					Amount:       new(big.Int).SetBytes(tronLog.Data).String(),
				})
			}
		}
	}
	if len(txRes.InternalTransactions) > 0 {
		for _, interTx := range txRes.InternalTransactions {
			if !interTx.Rejected {
				for _, callValueInfo := range interTx.CallValueInfo {
					if callValueInfo.CallValue != 0 && (callValueInfo.TokenId == "" || callValueInfo.TokenId == "trx") {
						transfers = append(transfers, &model.Transfer{
							ChainIndex: t.chainIndex,
							From:       encodeTronAddrWithChecksum(interTx.CallerAddress),
							To:         encodeTronAddrWithChecksum(interTx.TransferToAddress),
							Amount:     strconv.FormatInt(callValueInfo.CallValue, 10),
						})
					}
				}
			}

		}
	}
	return transfers, nil
}

func encodeTronAddr(address []byte) string {
	// Step 1: 添加前缀 0x41
	prefix := []byte{0x41}
	addressBytes := append(prefix, address...)

	return encodeTronAddrWithChecksum(addressBytes)
}

func encodeTronAddrWithChecksum(addressBytes []byte) string {
	// Step 2: 计算两次 SHA256 校验和
	firstSHA := sha256.Sum256(addressBytes)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]

	// Step 3: 拼接地址和校验和
	fullBytes := append(addressBytes, checksum...)

	// Step 4: Base58 编码
	return base58.Encode(fullBytes)
}

// ----------------- quicknode grpc token --------------------
type auth struct {
	token string
}

func (a *auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"x-token": a.token,
	}, nil
}

func (*auth) RequireTransportSecurity() bool {
	return false
}
