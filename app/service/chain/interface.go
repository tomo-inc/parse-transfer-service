package chain

import "github.com/tomo.inc/parse-transfer-service/app/model"

type Parser interface {
	FetchTxTransfer(tx string) ([]*model.Transfer, error)
}
