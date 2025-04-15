package model

type Transfer struct {
	ChainIndex   string `json:"chainIndex"`
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       string `json:"amount"`
	TokenAddress string `json:"tokenAddress"`
}
