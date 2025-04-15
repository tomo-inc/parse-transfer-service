package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tomo.inc/parse-transfer-service/app/handler/middleware"
	"github.com/tomo.inc/parse-transfer-service/app/service"
	"net/http"
)

type Handle struct {
	chainSvc service.ChainService
}

func NewHandle(chainSvc service.ChainService) *Handle {
	return &Handle{
		chainSvc: chainSvc,
	}
}

func (h *Handle) Router(engin *gin.Engine) {
	// middleware
	engin.Use(middleware.Logger())
	engin.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := engin.Group("/v1")
	v1.GET("/tx_transfers/:chainId/:txHash", h.FetchTransfer)
}
