package chain

import (
	"bytes"
	"crypto/tls"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"math/big"
	"strings"
	"testing"
)

func TestTron_getTransaction(t *testing.T) {
	client := client.NewGrpcClient("evocative-frosty-hill.tron-mainnet.quiknode.pro:50051")
	defer client.Stop()
	err := client.Start(grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithPerRPCCredentials(&auth{
			token: "ca682c8bf1e9e7c8981dd84da58c96eefe998882",
		}))
	require.NoError(t, err)
	res, err := client.GetTransactionInfoByID("8b9ad78916dc9b3384a8c6ef2d020f11fcd8e179a59dff7d8bc40f3564eae2f8")
	require.NoError(t, err)
	//t.Log(res.Log)
	//t.Log(res.InternalTransactions)

	for i, tronLog := range res.Log {
		if len(tronLog.Topics) == 3 && len(tronLog.Data) == 32 && bytes.Equal(tronLog.Topics[0], common.Hex2Bytes("ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")) {
			t.Logf("index: %d", i)
			t.Logf("contract: %s", encodeTronAddr(tronLog.Address))
			t.Logf("from: %s", encodeTronAddr(tronLog.Topics[1][12:]))
			t.Logf("to: %s", encodeTronAddr(tronLog.Topics[2][12:]))
			t.Logf("value: %s", new(big.Int).SetBytes(tronLog.Data).String())
		}
	}

	for i, interTx := range res.InternalTransactions {
		t.Logf("index: %d", i)
		t.Logf("caller: %s", encodeTronAddrWithChecksum(interTx.CallerAddress))
		t.Logf("callTo: %s", encodeTronAddrWithChecksum(interTx.TransferToAddress))
		for _, callValueInfo := range interTx.CallValueInfo {
			if callValueInfo.TokenId == "" || callValueInfo.TokenId == "trx" {
				t.Logf("value: %d", callValueInfo.CallValue)
			}
		}
		t.Logf("Rejected: %v", interTx.Rejected)
	}
}

func TestTron_FetchTxTransfer(t *testing.T) {
	tron := NewTron("evocative-frosty-hill.tron-mainnet.quiknode.pro:50051", "ca682c8bf1e9e7c8981dd84da58c96eefe998882", "1948400")
	data, err := tron.FetchTxTransfer("8b9ad78916dc9b3384a8c6ef2d020f11fcd8e179a59dff7d8bc40f3564eae2f8")
	require.NoError(t, err)

	for _, info := range data {
		t.Log(info)
	}
}

func Test_trim(t *testing.T) {
	t.Log(strings.Trim("0x00111", "0x"))
	require.Equal(t, strings.TrimPrefix("0x0111", "0x"), "0111")
	t.Log(strings.TrimLeft("0x0111", "0x"))
	t.Log(strings.TrimPrefix("0x0111", "0x"))
}
