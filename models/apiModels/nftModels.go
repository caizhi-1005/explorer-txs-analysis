package apiModels

//type RespContractList struct {
//	ContractAddress string `json:"contract_address"`
//	Name            string `json:"name"`
//	Symbol          string `json:"symbol"`
//	Logo            string `json:"logo"`
//}

type RespNftList struct {
	ContractAddress string `json:"contract_address"`
	Symbol          string `json:"symbol"`
	Logo            string `json:"logo"`
}

type ReqTokenIDList struct {
	ContractAddress string `json:"contract_address"`
	ReqCommon
}

type ReqNFTAddressDetail struct {
	ContractAddress string `json:"contract_address"`
	Address         string `json:"address"`
}

type RespNFTAddressDetail struct {
	AddressType         int                  `json:"address_type"`
	HoldTokenIds        []HoldTokenId        `json:"hold_token_ids"`
	HoldTokenIdsHistory []HoldTokenIdHistory `json:"hold_token_ids_history"`
	LongestHoldTokenId  string               `json:"longest_hold_token_id"`
	LongestHoldTime     string                `json:"longest_hold_time"`
	MaxProfitTokenId    string               `json:"max_profit_token_id"`
	MaxProfitValue      string               `json:"max_profit_value"`
}

type HoldTokenId struct {
	TokenId     string `json:"token_id"`
	AccountType int    `json:"account_type,omitempty"`
}

type HoldTokenIdHistory struct {
	HoldTokenIdHistory string `json:"hold_token_id_history"`
}

type HoldTokenIDList struct {
	HoldTokenId string `json:"hold_token_id"`
}

type ReqNFTTransferDetailsByAddress struct {
	ReqPagination
	ContractAddress string `json:"contract_address"`
	Address         string `json:"address"`
}

type RespNFTTransferDetailsByAddress struct {
	TokenId       string `json:"token_id"`
	TransferCount int    `json:"transfer_count"`
	From          string `json:"from"`
	To            string `json:"to"`
	FromType      int    `json:"from_type"`
	ToType        int    `json:"to_type"`
	Price         string `json:"price"`
}

type ReqNFTTxDetail struct {
	ContractAddress string `json:"contract_address"`
	TxHash          string `json:"tx_hash"`
}

type RespNFTTxDetail struct {
	TxHash        string `json:"tx_hash"`
	TokenId       string `json:"token_id"`
	TransferCount int    `json:"transfer_count"`
	Method        string `json:"method"`
	From          string `json:"from"`
	To            string `json:"to"`
	TxTime        string `json:"tx_time"`
}

type ReqNFTDetail struct {
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
}

type RespNFTDetail struct {
	Name               string `json:"name"`
	Symbol             string `json:"symbol"`
	Logo               string `json:"logo"`
	TokenType          string `json:"token_type"`
	Holder             string `json:"holder"`
	ContractAddress    string `json:"contract_address"`
	TransferCount      int    `json:"transfer_count"`
	HistoryHolderCount int    `json:"history_holder_count"`
	MintTime           string `json:"mint_time"`
	LongestHoldTime    string  `json:"longest_hold_time"`
	CurrentHoldTime    string `json:"current_hold_time"`
	LatestPrice        string `json:"latest_price"`
	HighestPrice       string `json:"highest_price"`
	LowestPrice        string `json:"lowest_price"`
}

type ReqNFTTransferDetailsByTokenId struct {
	ReqPagination
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
}

type RespNFTTransferDetailsByTokenId struct {
	TxTime   string `json:"tx_time"`
	Method   string `json:"method"`
	From     string `json:"from"`
	To       string `json:"to"`
	FromType int    `json:"from_type"`
	ToType   int    `json:"to_type"`
	Price    string `json:"price"`
}

type ReqNFTStartAnalysis struct {
	ReqCommon
	ContractAddress string `json:"contract_address"`
	TokenId         string `json:"token_id"`
}

type RespNFTStartAnalysis struct {
	TxTime   string `json:"tx_time"`
	TxType   string `json:"tx_type"`
	From     string `json:"from"`
	To       string `json:"to"`
	FromType int    `json:"from_type"`
	ToType   int    `json:"to_type"`
}

type ReqNFTTrace struct {
	ContractAddress string `json:"contract_address"`
	Address         string `json:"address"`
	TokenID         string `json:"token_id"`
	Count           string `json:"count"`
	Direction       string `json:"direction"` //in out
}

type RespNFTTrace struct {
	TxTime   string `json:"tx_time"`
	TxType   string `json:"tx_type"`
	From     string `json:"from"`
	To       string `json:"to"`
	FromType int    `json:"from_type"`
	ToType   int    `json:"to_type"`
}
