package chain

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	apperr "github.com/tomo.inc/parse-transfer-service/app/err"
	"github.com/tomo.inc/parse-transfer-service/app/model"
	"github.com/tomo.inc/parse-transfer-service/pkg/bot"
	"strconv"
)

type Solana struct {
	chainIndex string
	client     *rpc.Client
}

func NewSolana(endpoint string, chainIndex string) *Solana {
	client := rpc.New(endpoint)
	return &Solana{
		chainIndex: chainIndex,
		client:     client,
	}
}

func (s *Solana) FetchTxTransfer(txHash string) ([]*model.Transfer, error) {
	var ctx = context.Background()
	sig, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		return nil, apperr.InvalidTxErr
	}

	tx, err := s.client.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
		Commitment:                     rpc.CommitmentConfirmed,
		MaxSupportedTransactionVersion: lo.ToPtr(uint64(0)),
	})
	if err != nil {
		bot.SendMsg(context.Background(), bot.Msg{
			Title: "get transaction failed",
			Level: bot.Error,
			Data: []*bot.KeyValue{
				{
					Key:   "chain index",
					Value: s.chainIndex,
				},
				{
					Key:   "tx hash",
					Value: txHash,
				},
				{
					Key:   "err",
					Value: err.Error(),
				},
			},
		})
		log.Error().Err(err).Msg("Failed to get transaction")
		return nil, err
	}

	txParsed, err := tx.Transaction.GetTransaction()
	if err != nil {
		bot.SendMsg(context.Background(), bot.Msg{
			Title: "parse transaction failed",
			Level: bot.Error,
			Data: []*bot.KeyValue{
				{
					Key:   "chain index",
					Value: s.chainIndex,
				},
				{
					Key:   "tx hash",
					Value: txHash,
				},
				{
					Key:   "err",
					Value: err.Error(),
				},
			},
		})
		log.Error().Err(err).Msg("Failed to decode transaction")
		return nil, err
	}

	accountKeys := s.getAccountKeys(txParsed, tx)
	allInstruction := s.makeInstruction(txParsed, tx)
	tokenAccount := s.getAccount(tx.Meta.PreTokenBalances, tx.Meta.PostTokenBalances)
	return s.parseInstructions(accountKeys, tokenAccount, allInstruction), nil
}

func (s *Solana) getAccountKeys(txParsed *solana.Transaction, tx *rpc.GetTransactionResult) []solana.PK {
	accounts := make([]solana.PK, 0, len(txParsed.Message.AccountKeys)+len(tx.Meta.LoadedAddresses.Writable)+len(tx.Meta.LoadedAddresses.ReadOnly))
	accounts = append(accounts, txParsed.Message.AccountKeys...)
	accounts = append(accounts, tx.Meta.LoadedAddresses.Writable...)
	accounts = append(accounts, tx.Meta.LoadedAddresses.ReadOnly...)
	return accounts
}

type Instruction struct {
	Instruction       solana.CompiledInstruction
	InnerInstructions []solana.CompiledInstruction
}

func (s *Solana) makeInstruction(txParsed *solana.Transaction, tx *rpc.GetTransactionResult) []*Instruction {
	outers := make([]*Instruction, len(txParsed.Message.Instructions))
	for i, inst := range txParsed.Message.Instructions {
		outers[i] = &Instruction{Instruction: inst}
	}

	// 这里需要关注一下stack height 的问题, 不同的深度, 关于account 的信息会不同
	for _, innerInstuction := range tx.Meta.InnerInstructions {
		outers[innerInstuction.Index].InnerInstructions = innerInstuction.Instructions
	}
	return outers
}

type account struct {
	Owner solana.PublicKey
	Mint  solana.PublicKey
}

func (s *Solana) getAccount(preTokenBalance, postTokenBalance []rpc.TokenBalance) map[uint16]*account {
	tokenAccounts := make(map[uint16]*account)
	for _, tokenBalance := range preTokenBalance {
		tokenAccounts[tokenBalance.AccountIndex] = &account{
			Owner: lo.FromPtr(tokenBalance.Owner),
			Mint:  tokenBalance.Mint,
		}
	}
	for _, tokenBalance := range postTokenBalance {
		tokenAccounts[tokenBalance.AccountIndex] = &account{
			Owner: lo.FromPtr(tokenBalance.Owner),
			Mint:  tokenBalance.Mint,
		}
	}
	return tokenAccounts
}

func (s *Solana) parseInstructions(accounts []solana.PK, tokenAccounts map[uint16]*account, instructions []*Instruction) (transfers []*model.Transfer) {
	transfers = make([]*model.Transfer, 0, len(instructions))
	for _, instruction := range instructions {
		transfer := s.parseInstruction(accounts, tokenAccounts, &instruction.Instruction)
		if transfer != nil {
			transfers = append(transfers, transfer)
		}
		for _, innerInstructions := range instruction.InnerInstructions {
			transfer := s.parseInstruction(accounts, tokenAccounts, &innerInstructions)
			if transfer != nil {
				transfers = append(transfers, transfer)
			}
		}
	}
	return transfers
}

func (s *Solana) parseInstruction(accounts []solana.PK, tokenAccounts map[uint16]*account, instruction *solana.CompiledInstruction) *model.Transfer {
	log.Debug().Uint16("progromID", instruction.ProgramIDIndex).Int("Accounts len", len(instruction.Accounts)).Int("data len", len(instruction.Data)).Str("data", hex.EncodeToString(instruction.Data)).Msg("show")

	switch accounts[instruction.ProgramIDIndex] {
	case system.ProgramID:
		if len(instruction.Accounts) == 2 && len(instruction.Data) == 12 {
			amount := binary.LittleEndian.Uint64(instruction.Data[4:])
			return &model.Transfer{
				ChainIndex:   s.chainIndex,
				From:         accounts[instruction.Accounts[0]].String(),
				To:           accounts[instruction.Accounts[1]].String(),
				Amount:       strconv.FormatUint(amount, 10),
				TokenAddress: SOL_TOKEN_ADDRESS,
			}
		}
	case solana.TokenProgramID, solana.Token2022ProgramID:
		if len(instruction.Data) >= 1 {
			switch instruction.Data[0] {
			case byte(1):
				log.Info().Any("account", tokenAccounts[instruction.Accounts[0]]).Msg("show account")
				//token.InitializeAccount{}
				// https://solscan.io/tx/4ReDF9XgzyBBc8XU6m8urFhy7oCU7UeuwVTniJHqukhPFfMWRhL9ZdyJzWHMw1YTFKAHQwYVHxdrabR2VKmM5gjf
				if len(instruction.Accounts) == 4 && len(instruction.Data) == 1 {
					if tokenAccounts[instruction.Accounts[0]] == nil {
						tokenAccounts[instruction.Accounts[0]] = &account{
							Owner: accounts[instruction.Accounts[2]],
							Mint:  accounts[instruction.Accounts[1]],
						}
						log.Info().Any("account", tokenAccounts[instruction.Accounts[0]]).Msg("show account")
					}
				}
			case byte(16):
				//token.InitializeAccount2{}
				if len(instruction.Accounts) == 3 && len(instruction.Data) == 33 {
					initializeAccount2 := token.InitializeAccount2{}
					err := bin.NewBorshDecoder(instruction.Data[1:]).Decode(&initializeAccount2)
					if err != nil {
						return nil
					}
					if tokenAccounts[instruction.Accounts[0]] == nil {
						tokenAccounts[instruction.Accounts[0]] = &account{
							Owner: lo.FromPtr(initializeAccount2.Owner),
							Mint:  accounts[instruction.Accounts[1]],
						}
					}
				}
			case byte(18):
				//token.InitializeAccount3{}
				// https://solscan.io/tx/45SX1KP9nYyh4S3rumUYgjLweSN3kuy9HKRR3Bb3xBkFHqDwcY6KGBG4RQ1xzGE2XWT3HbcAbNV3bTdB5ELMZuUw
				if len(instruction.Accounts) == 2 && len(instruction.Data) == 33 {
					initializeAccount3 := token.InitializeAccount3{}
					err := bin.NewBorshDecoder(instruction.Data[1:]).Decode(&initializeAccount3)
					if err != nil {
						return nil
					}
					if tokenAccounts[instruction.Accounts[0]] == nil {
						tokenAccounts[instruction.Accounts[0]] = &account{
							Owner: lo.FromPtr(initializeAccount3.Owner),
							Mint:  accounts[instruction.Accounts[1]],
						}
					}
				}

			case byte(3):
				//token.Transfer{}
				if len(instruction.Accounts) == 3 && len(instruction.Data) == 9 {
					amount := binary.LittleEndian.Uint64(instruction.Data[1:])

					if tokenAccounts[instruction.Accounts[1]] == nil {
						// 账户有可能会被 close 掉
						// for example: https://solscan.io/tx/4ReDF9XgzyBBc8XU6m8urFhy7oCU7UeuwVTniJHqukhPFfMWRhL9ZdyJzWHMw1YTFKAHQwYVHxdrabR2VKmM5gjf
						return nil
					}
					return &model.Transfer{
						ChainIndex:   s.chainIndex,
						From:         accounts[instruction.Accounts[2]].String(),            // owner (from)
						To:           tokenAccounts[instruction.Accounts[1]].Owner.String(), // destination.owner (to)
						Amount:       strconv.FormatUint(amount, 10),
						TokenAddress: tokenAccounts[instruction.Accounts[1]].Mint.String(),
					}
				}
			case byte(12):
				if len(instruction.Accounts) == 4 && len(instruction.Data) == 10 {
					transferChecked := token.TransferChecked{}
					err := bin.NewBorshDecoder(instruction.Data[1:]).Decode(&transferChecked)
					if err != nil {
						return nil
					}
					if tokenAccounts[instruction.Accounts[2]] == nil {
						// 账户有可能会被 close 掉
						// for example: https://solscan.io/tx/4ReDF9XgzyBBc8XU6m8urFhy7oCU7UeuwVTniJHqukhPFfMWRhL9ZdyJzWHMw1YTFKAHQwYVHxdrabR2VKmM5gjf
						return nil
					}

					return &model.Transfer{
						ChainIndex:   s.chainIndex,
						From:         accounts[instruction.Accounts[3]].String(),
						To:           tokenAccounts[instruction.Accounts[2]].Owner.String(),
						Amount:       strconv.FormatUint(lo.FromPtr(transferChecked.Amount), 10),
						TokenAddress: accounts[instruction.Accounts[1]].String(),
					}
				}
			}
		}
	}
	return nil
}
