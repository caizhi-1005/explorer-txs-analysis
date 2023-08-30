package nebulaModels

import (
	"encoding/json"
)

type TransactionEdge struct {
	TxHash      string `norm:"tx_hash"`
	TxTime      string `norm:"tx_time"`
	FromAddress string `norm:"from_address"`
	ToAddress   string `norm:"to_address"`
	Amount      string `norm:"amount"`
}

type RouteTxStep struct {
	Transaction []TransactionEdge `json:"transaction"`
	Src         string            `json:"from,omitempty"`
	Dst         string            `json:"to,omitempty"`
}

type TxsRoute struct {
	Steps []RouteTxStep `json:"tx_steps"`
}

func (r TxsRoute) String() string {
	d, _ := json.Marshal(r)
	return string(d)
}

//-------------------------------

type TxsEdge struct {
	BlockId     int     `norm:"block_id"`
	TxHash      string  `norm:"tx_hash"`
	TxTime      string  `norm:"tx_time"`
	FromAddress string  `norm:"from_address"`
	ToAddress   string  `norm:"to_address"`
	Value       string  `norm:"value"`
	Amount      float64 `norm:"amount"`
	//TokenId      string  `norm:"token_id"`
	TokenId      string `json:"token_id,omitempty"`
	TxType       int    `norm:"tx_type"`
	Caller       string `norm:"caller"`
	Callee       string `norm:"callee"`
	TokenAddress string `norm:"token_address"`
}

type TxStep struct {
	Txs []TxsEdge `json:"txs"`
	Src Address   `json:"from_address,omitempty"`
	Dst Address   `json:"to_address,omitempty"`
}

type AddressTag struct {
	Type    int    `json:"type"`
	Address string `json:"address"`
}

type TxSteps struct {
	Steps []TxStep `json:"tx_steps"`
}

func (r TxSteps) String() string {
	d, _ := json.Marshal(r)
	return string(d)
}
