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
	ContractAddress string `json:"contract_address"`
	Address         string `json:"address"`
	Count           string `json:"count"`     //展开数量
	Direction       string `json:"direction"` //in out all
}

type RespAddressTxGraph struct {
	//当前地址，地址类型
	//目标地址，地址类型
	//交易次数
	//交易金额
	Count     string `json:"count"`
	Type      string `json:"type"` //0:全部， -1:前前， 1:向后
	Address   string `json:"address"`
	Direction string `json:"direction"` //in out
}

type TxEdge struct {
	From                 string  `json:"from,omitempty"`
	SrcAddress           string  `json:"src_address,omitempty"`
	Index                int     `json:"index,omitempty"`
	FlowType             int     `json:"flow_type,omitempty"`
	To                   string  `json:"to,omitempty"`
	TotalValue           float64 `json:"total_value,omitempty"`
	TxCount              int     `json:"tx_count,omitempty"`
	FirstTransactionTime int     `json:"first_transaction_time,omitempty"`
	LastTransactionTime  int     `json:"last_transaction_time,omitempty"`
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

type RespAddressTxList struct {
	TxTime string `json:"tx_time"`
	TxHash string `json:"tx_hash"`
	Amount string `json:"amount"`
	Symbol string `json:"symbol,omitempty"`
}
