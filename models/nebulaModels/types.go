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
