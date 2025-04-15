package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tomo.inc/parse-transfer-service/app/handler/response"
)

func (h *Handle) FetchTransfer(c *gin.Context) {
	chainIndex := c.Param("chainId")
	txhash := c.Param("txHash")
	data, err := h.chainSvc.FetchTxTransfer(chainIndex, txhash)
	if err != nil {
		response.HandlerError(c, err)
	} else {
		response.Success(c, data)
	}
	return
}
