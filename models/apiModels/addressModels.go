package apiModels

type RespAddressDetail struct {
	Address            string  `json:"address"`
	Type               int     `json:"type"`
	OutAddressCount    int     `json:"out_address_count"`
	InAddressCount     int     `json:"in_address_count"`
	Balance            float64 `json:"balance"`
	FirstTxTime        string  `json:"first_tx_time"`
	TxCount            int     `json:"tx_count"`
	MaxTxAmount        float64 `json:"max_tx_amount"`
	ReceiveAmountTotal float64 `json:"receive_amount_total"`
	SendAmountTotal    float64 `json:"send_amount_total"`
}

type ReqTxAnalysis struct {
	Address string `json:"address"`
	Type    string `json:"type"` // 1:CMP交易  2:CRC20交易
}

type RespAddressTxAnalysis struct {
	AddressVertex []AddressVertex `json:"nodes"`
	Edges         []TxEdge        `json:"edges"`
}

type AddressVertex struct {
	Address    string `json:"address,omitempty"`
	SrcAddress string `json:"srcAddress,omitempty"`
	IsInit     bool   `json:"isInit,omitempty"`
}

type TxEdge struct {
	From                 string  `json:"from,omitempty"`
	SrcAddress           string  `json:"srcAddress,omitempty"`
	Index                int     `json:"index,omitempty"`
	FlowType             int     `json:"flowType,omitempty"`
	To                   string  `json:"to,omitempty"`
	TotalValue           float64 `json:"totalValue,omitempty"`
	TxCount              int     `json:"txCount,omitempty"`
	FirstTransactionTime int     `json:"firstTransactionTime,omitempty"`
	LastTransactionTime  int     `json:"lastTransactionTime,omitempty"`
}

type Nodes struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type Edges struct {
	Source    string  `json:"source"`
	Target    string  `json:"target"`
	TxCount   float64 `json:"tx_count"`
	UsdtCount float64 `json:"usdt_count"`
}

type RespTxDetailInfo struct {
	From         string  `json:"from"`
	FromType     int     `json:"from_type"`
	To           string  `json:"to"`
	ToType       int     `json:"to_type"`
	TxCount      float64 `json:"tx_count"`
	TxAmount     float64 `json:"tx_amount"`
	FirstTxTime  string  `json:"first_tx_time"`
	LatestTxTime string  `json:"latest_tx_time"`
}

type ReqAddressTxDetail struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ReqAddressTxList struct {
	ReqAddressTxDetail
	ReqPagination
}
