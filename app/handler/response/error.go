package response

import (
	"github.com/gin-gonic/gin"
	apperr "github.com/tomo.inc/parse-transfer-service/app/err"
	"net/http"
)

func HandlerError(c *gin.Context, err error) {
	var res = Response[interface{}]{}
	switch err.(type) {
	case *apperr.CustomError:
		cusErr := err.(*apperr.CustomError)
		res.Code = cusErr.Code
		res.Message = cusErr.Message
		c.JSON(cusErr.Status, res)
	default:
		res.Code = apperr.CodeErr
		res.Message = "Internal Server Error"
		c.JSON(http.StatusInternalServerError, res)
	}
	return
}
