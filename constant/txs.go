package constant

const (
	BlockId      = "block_id"
	TxHash       = "tx_hash"
	TxTime       = "tx_time"
	FromAddress  = "from_address"
	ToAddress    = "to_address"
	Value        = "value"
	Amount       = "amount"
	TokenId      = "token_id"
	TxType       = "tx_type"
	TokenAddress = "token_address"
)

const (
	TxRanK   = "Tx_Rank"
	CoinTxs  = "Coin_Txs"
	TokenTxs = "Token_Txs"
	NftTxs   = "NFT_Txs"
)

const (
	DirPath    = "rankfile/"
	RankSuffix = ".rank"
)

const (
	TransferCode      = "0xa9059cbb"
	TransferEventCode = "0xddf252ad"
	MintCode          = "0xa0712d68"

	TransferMethod = "Transfer"
	MintMethod     = "Mint"
)

const BASE_TOKEN_DECIMAL float64 = 18
