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
	EVMEndpoints map[string]EvmConfig,
	SOlEndpoints map[string]string,
	TRONEndpoints map[string]string,
) ChainService {
	chains := make(map[string]chain.Parser)
	for chainIndex, config := range EVMEndpoints {
		chains[chainIndex] = chain.NewEvm(config.SupportDebug, config.Endpoint, chainIndex)
	}
	for chainIndex, endpoint := range SOlEndpoints {
		chains[chainIndex] = chain.NewSolana(endpoint, chainIndex)
	}
	for chainIndex, endpoint := range TRONEndpoints {
		chains[chainIndex] = chain.NewTron(endpoint, chainIndex)
	}

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
