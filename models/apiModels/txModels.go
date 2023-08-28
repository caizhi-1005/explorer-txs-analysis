package apiModels

type RespTxDetail struct {
	TxHash   string  `json:"tx_hash"`
	Amount   float64 `json:"amount"`
	TxFee    float64 `json:"tx_fee"`
	TxTime   string  `json:"tx_time"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	FromType int     `json:"from_type"`
	ToType   int     `json:"to_type"`
}

type ReqAddressTxGraph struct {
	Count   string `json:"count"`
	Type    string `json:"type"` //0:全部， -1:前前， 1:向后
	Address string `json:"address"`
}

type RespTxAddressDetail struct {
	Address            string  `json:"address"`
	Type               int     `json:"type"`
	OutAddressCount    int     `json:"out_address_count"`
	InAddressCount     int     `json:"in_address_count"`
	Balance            float64 `json:"balance"`
	TxCount            int     `json:"tx_count"`
	ReceiveAmountTotal float64 `json:"receive_amount_total"`
	SendAmountTotal    float64 `json:"send_amount_total"`
	InTxCount          int     `json:"in_tx_count"`
	OutTxCount         int     `json:"out_tx_count"`
}
