package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	apperr "github.com/tomo.inc/parse-transfer-service/app/err"
	"testing"
)

var (
	chainIndex    = "5600"
	testEvmClient = NewEvm(true, "https://evocative-frosty-hill.bsc.quiknode.pro/ca682c8bf1e9e7c8981dd84da58c96eefe998882", chainIndex)
)

func Test_evm_parseLog(t *testing.T) {
	transfer, err := testEvmClient.parseLog(common.HexToHash("0x658fc1c51da8d319e91affbb7dfe10d099817e11e599a460a1c296edb7c1c330"))
	require.NoError(t, err)
	require.Len(t, transfer, 2)
	require.True(t, transfer[0].ChainIndex == chainIndex)
	require.True(t, transfer[0].From == "0xb1000058c87D843Fc0154591Ff9d72AF5e7213D5")
	require.True(t, transfer[0].To == "0xF87E6FC700fa42963dfB57010fbC8f902b2c0fC2")
	require.True(t, transfer[0].Amount == "99294533141400000")
	require.True(t, transfer[0].TokenAddress == "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")
}

func Test_evm_parseDebug(t *testing.T) {
	transfers, err := testEvmClient.parseDebug(common.HexToHash("0x160bda5500fff93777d965ac1c8083957651d3c3b7b660ef30fa3ff8c624231f"))
	require.NoError(t, err)
	require.Len(t, transfers, 2)
	require.True(t, transfers[1].ChainIndex == chainIndex)
	require.True(t, transfers[1].From == "0x1a0a18ac4becddbd6389559687d1a73d8927e416")
	require.True(t, transfers[1].To == "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c")
	require.True(t, transfers[1].Amount == "10000000000000000000")
}

func Test_evm_FetchTxTransfer(t *testing.T) {
	_, err := testEvmClient.FetchTxTransfer("0x160bda5500fff93777d965ac1c8083957651d3c3b7b660ef30fa3ff8c62")
	require.ErrorIs(t, err, apperr.InvalidTxErr)

	transfers, err := testEvmClient.FetchTxTransfer("0x160bda5500fff93777d965ac1c8083957651d3c3b7b660ef30fa3ff8c624231f")
	require.NoError(t, err)
	require.Len(t, transfers, 6)
}
