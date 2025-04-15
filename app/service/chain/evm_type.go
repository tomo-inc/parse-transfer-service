package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Config struct {
	EnableMemory     bool   `json:"enableMemory"`     // enable memory capture
	DisableStack     bool   `json:"disableStack"`     // disable stack capture
	DisableStorage   bool   `json:"disableStorage"`   // disable storage capture
	EnableReturnData bool   `json:"enableReturnData"` // enable return data capture
	Tracer           string `json:"tracer,omitempty"` //
	Timeout          string `json:"timeout,omitempty"`
}

type CallFrame struct {
	Type         string       `json:"type"`
	From         string       `json:"from"`
	To           string       `json:"to,omitempty"`
	Value        string       `json:"value,omitempty"`
	Gas          string       `json:"gas"`
	GasUsed      string       `json:"gasUsed"`
	Input        string       `json:"input"`
	Output       string       `json:"output,omitempty"`
	Error        string       `json:"error,omitempty"`
	Calls        []*CallFrame `json:"calls,omitempty"`
	RevertReason string       `json:"revertReason,omitempty"`
	Logs         []callLog    `json:"logs,omitempty" rlp:"optional"`
}

type callLog struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    hexutil.Bytes  `json:"data"`
}
