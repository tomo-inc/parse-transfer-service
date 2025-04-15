package err

import "net/http"

var (
	CodeErr = -1
)

var (
	InternalErr          = NewCustomError(http.StatusInternalServerError, CodeErr, "internal error")
	InvalidChainIndexErr = NewCustomError(http.StatusBadRequest, CodeErr, "invalid chain index")
	NotSupportChainErr   = NewCustomError(http.StatusNotAcceptable, CodeErr, "not support chain")
	InvalidTxErr         = NewCustomError(http.StatusBadRequest, CodeErr, "invalid tx error")
)
