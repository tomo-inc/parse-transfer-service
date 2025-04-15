package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var CodeOk = 1

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func Success[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Response[T]{
		Code: CodeOk,
		Data: data,
	})
}

type HealthyResponse struct {
	Healthy bool `json:"healthy"` // health status for the rpc service
}
