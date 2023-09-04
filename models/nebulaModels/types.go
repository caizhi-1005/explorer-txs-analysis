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

//type TxsEdge struct {
//	BlockId      int64   `norm:"block_id"`
//	TxHash       string  `norm:"tx_hash"`
//	TxTime       string  `norm:"tx_time"`
//	FromAddress  string  `norm:"from_address"`
//	ToAddress    string  `norm:"to_address"`
//	Value        string  `norm:"value"`
//	Amount       float64 `norm:"amount"`
//	TokenId      string  `norm:"token_id"`
//	TxType       int64   `norm:"tx_type"`
//	Caller       string  `norm:"caller"`
//	Callee       string  `norm:"callee"`
//	TokenAddress string  `norm:"token_address"`
//}

type TxsEdge struct {
	BlockId      int64   `json:"block_id"`
	TxHash       string  `json:"tx_hash"`
	TxTime       string  `json:"tx_time"`
	FromAddress  string  `json:"from_address"`
	ToAddress    string  `json:"to_address"`
	Value        string  `json:"value"`
	Amount       float64 `json:"amount"`
	TokenId      string  `json:"token_id"`
	TxType       int64   `json:"tx_type"`
	Caller       string  `json:"caller,omitempty"`
	Callee       string  `json:"callee,omitempty"`
	TokenAddress string  `json:"token_address"`
	TotalAmount  float64 `json:"total_amount,omitempty"`
	TxCount      int64   `json:"tx_count,omitempty"`
}

type TxStep struct {
	Txs []TxsEdge `json:"txs"`
	Src Address   `json:"from_address,omitempty"`
	Dst Address   `json:"to_address,omitempty"`
}

type AddressTag struct {
	Type       int    `json:"type"` //账户类型：1-account 2-contract 0-无效值
	Address    string `json:"address"`
	SrcAddress string `json:"src_address"`
}

type TxSteps struct {
	Steps []TxStep `json:"tx_steps"`
}

//func (r TxSteps) String() string {
//	d, _ := json.Marshal(r)
//	return string(d)
//}

type RespGraph struct {
	Edge   []*TxsEdge    `json:"edge"`
	Vertex []*AddressTag `json:"vertex"`
}
