package apiModels

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
	AddressType            int                   `json:"address_type"`
	HoldTokenIdList        []HoldTokenIds        `json:"hold_token_id_list"`
	HoldTokenIdHistoryList []HoldTokenIdsHistory `json:"hold_token_id_history_list"`
	LongestHoldTokenId     string                `json:"longest_hold_token_id"`
	LongestHoldTime        string                `json:"longest_hold_time"`
}

type HoldTokenIds struct {
	HoldTokenId string `json:"token_id"`
	AccountType int    `json:"account_type"`
}

type HoldTokenIdsHistory struct {
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
	From          string `json:"from"`
	To            string `json:"to"`
	TransferCount int    `json:"transfer_count"`
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
	TokenType          string `json:"token_type"`
	Holder             string `json:"holder"`
	ContractAddress    string `json:"contract_address"`
	TransferCount      int    `json:"transfer_count"`
	HistoryHolderCount int    `json:"history_holder_count"`
	//MintTime           time.Time `json:"mint_time,omitempty" copier:"-"`
	MintTime        string `json:"mint_time"`
	LongestHoldTime string `json:"longest_hold_time"`
	CurrentHoldTime string `json:"current_hold_time"`
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
}

type ReqNFTStartAnalysis struct {
	ReqCommon
	ContractAddress string `json:"contract_address"`
	Input           string `json:"input"`
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
	Address         string `json:"_address"`
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
