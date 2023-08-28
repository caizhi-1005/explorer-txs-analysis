package nebulaModels

import (
	"github.com/zhihu/norm"
	"time"
)

type Address struct {
	norm.VModel
	Address string `norm:"address"`
	Type    int    `norm:"type"`
}

type CoinTxs struct {
	norm.EModel
	BlockId     int64     `norm:"block_id"`
	TxHash      string    `norm:"tx_hash"`
	TxTime      time.Time `norm:"tx_time"`
	FromAddress string    `norm:"from_address"`
	ToAddress   string    `norm:"to_address"`
	Value       string    `norm:"value"`
	Amount      float64   `norm:"amount"`
	Caller      string    `norm:"caller"`
	Callee      string    `norm:"callee"`
}

type TokenTxs struct {
	norm.EModel
	BlockId     int64     `norm:"block_id"`
	TxHash      string    `norm:"tx_hash"`
	TxTime      time.Time `norm:"tx_time"`
	FromAddress string    `norm:"from_address"`
	ToAddress   string    `norm:"to_address"`
	Amount      string    `norm:"amount"`
	Caller      string    `norm:"caller"`
	Callee      string    `norm:"callee"`
}

type NFTTxs struct {
	norm.EModel
	BlockId     int64     `norm:"block_id"`
	TxHash      string    `norm:"tx_hash"`
	TxTime      time.Time `norm:"tx_time"`
	FromAddress string    `norm:"from_address"`
	ToAddress   string    `norm:"to_address"`
	Amount      string    `norm:"amount"`
	TokenId     string    `norm:"token_id"`
	TxType      int       `norm:"tx_type"`
	Caller      string    `norm:"caller"`
	Callee      string    `norm:"callee"`
}

type TxRank struct {
	norm.EModel
	Key    string `norm:"key"`
	TxRank int64  `norm:"tx_rank"`
}

var _ norm.IVertex = new(Address)
var _ norm.IEdge = new(CoinTxs)
var _ norm.IEdge = new(TokenTxs)
var _ norm.IEdge = new(NFTTxs)
var _ norm.IEdge = new(TxRank)

func (*Address) TagName() string {
	return "address"
}

func (t *Address) GetVid() interface{} {
	return t.Address
}

func (p *CoinTxs) EdgeName() string {
	return "Coin_Txs"
}

func (p *TokenTxs) EdgeName() string {
	return "Token_Txs"
}

func (p *NFTTxs) EdgeName() string {
	return "NFT_Txs"
}

func (p *TxRank) EdgeName() string {
	return "Tx_Rank"
}
