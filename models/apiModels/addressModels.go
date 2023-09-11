package apiModels

type ReqContractList struct {
	Name   string `json:"name"`
	TxType string `json:"tx_type"`
}

type RespContractList struct {
	Symbol string `json:"label"`
	//Symbol          string `json:"symbol"`
	//ContractAddress string `json:"contract_address"`
	ContractAddress string `json:"value"`
	Name            string `json:"name"`
	Logo            string `json:"logo"`
}

type ReqAddressDetail struct {
	ReqCommon
	ContractAddress string `json:"contract_address"`
}


type RespContractAddressInfo struct {
	Address     string  `json:"address"`
	AccountType int     `json:"account_type"`
	Balance     float64 `json:"balance"`
	Symbol      string  `json:"symbol,omitempty"`
	Decimals    int64   `json:"decimals,omitempty"`
}

type RespAddressDetail struct {
	Address            string  `json:"address"`
	Type               int     `json:"type"` // 账户类型：1-account 2-contract
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
	From            string `json:"from"`
	To              string `json:"to"`
	ContractAddress string `json:"contract_address"`
}

type ReqAddressTxList struct {
	ReqAddressTxDetail
	ReqPagination
}
