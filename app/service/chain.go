package service

import (
	apperr "github.com/tomo.inc/parse-transfer-service/app/err"
	"github.com/tomo.inc/parse-transfer-service/app/model"
	"github.com/tomo.inc/parse-transfer-service/app/service/chain"
)

type ChainService interface {
	FetchTxTransfer(chainIndex string, tx string) ([]*model.Transfer, error)
}

var _ ChainService = (*chainServiceImpl)(nil)

type chainServiceImpl struct {
	chains map[string]chain.Parser
}

func NewChain(
	chains map[string]chain.Parser,
) ChainService {
	return &chainServiceImpl{
		chains: chains,
	}
}

func (c *chainServiceImpl) FetchTxTransfer(chainIndex string, tx string) ([]*model.Transfer, error) {
	if res, ok := c.chains[chainIndex]; !ok {
		return nil, apperr.NotSupportChainErr
	} else {
		return res.FetchTxTransfer(tx)
	}
}
