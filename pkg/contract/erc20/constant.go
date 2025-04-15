package erc20

import "github.com/ethereum/go-ethereum/common"

var (
	TransferErc20Filter, _ = NewErc20Filterer(common.Address{}, nil)
)
