package chain

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
	"testing"
)

var testSolClient = NewSolana("https://practical-cosmological-meme.solana-mainnet.quiknode.pro/06261eba3f2e921d77444957ddb5c41140925813", "50100")

func Test_sol_FetchTxTransfer1(t *testing.T) {
	res, err := testSolClient.FetchTxTransfer("45SX1KP9nYyh4S3rumUYgjLweSN3kuy9HKRR3Bb3xBkFHqDwcY6KGBG4RQ1xzGE2XWT3HbcAbNV3bTdB5ELMZuUw")
	require.NoError(t, err)
	for _, r := range res {
		t.Log(r)
	}

}

//
//func Test_sol_PublicKey(t *testing.T) {
//	systemProgramID := solana.MPK("11111111111111111111111111111111")
//	t.Log(systemProgramID == system.ProgramID)
//}

func Test_sol_solTransfer(t *testing.T) {
	transferData, err := base58.Decode("3Bxs4Z6oyhaczjLK")
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(transferData), "02000000c0c62d0000000000")
	amount := binary.LittleEndian.Uint64(transferData[4:])
	require.True(t, amount == 3000000)
}

func Test_sol_tokenTransfer(t *testing.T) {
	transferData, err := base58.Decode("3GgDkdswmzXq")
	require.NoError(t, err)
	//t.Log(hex.EncodeToString(transferData))
	//require.Equal(t, hex.EncodeToString(transferData), "03133d975b15000000")
	amount := binary.LittleEndian.Uint64(transferData[1:])
	require.True(t, amount == 91730951443)
	//t.Log(amount)
}

//func Test_sol_token2022Transfer(t *testing.T) {
//	transferData, err := base58.Decode("3GgDkdswmzXq")
//	require.NoError(t, err)
//	//t.Log(hex.EncodeToString(transferData))
//	//require.Equal(t, hex.EncodeToString(transferData), "03133d975b15000000")
//	amount := binary.LittleEndian.Uint64(transferData[1:])
//	require.True(t, amount == 91730951443)
//	//t.Log(amount)
//}

func Test_sol_token_parseInstruction(t *testing.T) {
	client := new(Solana)

	accounts := []solana.PK{
		solana.TokenProgramID,
		solana.MPK("2tFEeixKGYjm8qsMVWs6QqXW6TNQrigCvwMpFDrwz67V"),
		solana.MPK("8ctcHN52LY21FEipCjr1MVWtoZa1irJQTPyAaTj72h7S"),
		solana.MPK("GtQZvXdBouuu3Yah3J2WGaVprLkvnYMWkbnhcY7ePhju"),
	}
	instructionData, _ := base58.Decode("3GgDkdswmzXq")
	instruction := &solana.CompiledInstruction{
		ProgramIDIndex: 0,
		Accounts:       []uint16{1, 2, 3}, // source destination owner
		Data:           instructionData,
	}
	tokenAccounts := map[uint16]*account{
		2: {
			Owner: solana.MPK("2tFEeixKGYjm8qsMVWs6QqXW6TNQrigCvwMpFDrwz67V"),
			Mint:  solana.MPK("So11111111111111111111111111111111111111112"),
		},
	}

	// from(owner) => to(destination.owner)
	transfer := client.parseInstruction(accounts, tokenAccounts, instruction)
	require.True(t, transfer.From == accounts[instruction.Accounts[2]].String())
	require.True(t, transfer.To == tokenAccounts[instruction.Accounts[1]].Owner.String())
	require.True(t, transfer.TokenAddress == tokenAccounts[instruction.Accounts[1]].Mint.String())
	require.True(t, transfer.Amount == "91730951443")
}

func Test_sol_token2022_parseInstruction(t *testing.T) {
	client := new(Solana)

	accounts := []solana.PK{
		solana.Token2022ProgramID,
		solana.MPK("EnVya7VPHFZhF7y4Mg95Y9TMD1v5Y2e8a27FGJkBanJH"),
		solana.MPK("E2rXoYHr8dxMUyYFZZaDVmgTiag29VnPkds3qdfSPump"),
		solana.MPK("C6Et8piGkumTQW9E2nN73hgZEUG8jyUSSakZyjKcmy8b"),
		solana.MPK("BCMiH3ZunQCmuHUxh2iN3zisRBvPd956c4gRWUNBprct"),
	}
	instructionData, _ := base58.Decode("ip2MSWDq5UKVo")
	instruction := &solana.CompiledInstruction{
		ProgramIDIndex: 0,
		Accounts:       []uint16{1, 2, 3, 4}, // source destination owner
		Data:           instructionData,
	}
	tokenAccounts := map[uint16]*account{
		3: {
			Owner: solana.MPK("B82qFxLHo9fN1Rj13Ncd6wGZBQpffciHFDFGuHSDm15Y"),
			Mint:  solana.MPK("E2rXoYHr8dxMUyYFZZaDVmgTiag29VnPkds3qdfSPump"),
		},
	}

	// from(owner) => to(destination.owner)
	transfer := client.parseInstruction(accounts, tokenAccounts, instruction)
	//t.Log(transfer)
	require.True(t, transfer.From == accounts[instruction.Accounts[3]].String())
	require.True(t, transfer.To == tokenAccounts[instruction.Accounts[2]].Owner.String())
	require.True(t, transfer.TokenAddress == accounts[instruction.Accounts[1]].String())
	require.True(t, transfer.Amount == "74381706196")
}
