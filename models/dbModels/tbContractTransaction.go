package dbModels

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/utils"
	"math/big"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type TbContractTransaction struct {
	Id              int64     `orm:"column(id);pk"`
	BlockId         int64     `orm:"column(block_id);null" description:"区块号"`
	BlockHash       string    `orm:"column(block_hash);size(255);null" description:"区块hash"`
	TxHash          string    `orm:"column(tx_hash);size(255);null" description:"交易hash"`
	TxTime          time.Time `orm:"column(tx_time);type(datetime);null" description:"交易时间"`
	ContractAddress string    `orm:"column(contract_address);size(255);null" description:"合约地址"`
	TokenAddress    string    `orm:"column(token_address);size(255);null" description:"token地址"`
	TokenType       int       `orm:"column(token_type);null" description:"token类型 1-erc20 2-erc721"`
	From            string    `orm:"column(from);size(255);null" description:"发起地址"`
	To              string    `orm:"column(to);size(255);null" description:"到大地址"`
	Value           string    `orm:"column(value);size(2048);null" description:"交易值"`
	Amount          float64   `orm:"column(amount);null;digits(64);decimals(18)" description:"交易金额"`
	TokenId         string    `orm:"column(token_id);size(2048);null" description:"tokenId"`
	TxStatus        int       `orm:"column(tx_status);null" description:"交易状态"`
	AddressList     string    `orm:"column(address_list);size(255);null" description:"地址列表(from|to|contract)"`
	Status          int       `orm:"column(status);null" description:"状态,1-正常 2-冻结"`
	IsDeleted       int8      `orm:"column(is_deleted);null" description:"删除状态 0-正常 1-删除"`
	SyncTime        time.Time `orm:"column(sync_time);type(timestamp);null;auto_now_add" description:"同步时间"`
	CreateTime      time.Time `orm:"column(create_time);type(timestamp);null;auto_now_add" description:"创建时间"`
	UpdateTime      time.Time `orm:"column(update_time);type(timestamp);null;auto_now_add" description:"更新时间"`
}

func (t *TbContractTransaction) TableName() string {
	return "tb_contract_transaction"
}

func init() {
	orm.RegisterModel(new(TbContractTransaction))
}

// NFTTransferDetails NFT溯源-NFT流转详情列表
func NFTTransferDetails(req apiModels.ReqNFTTransferDetailsByAddress) ([]*apiModels.RespNFTTransferDetailsByAddress, error) {
	list := make([]*apiModels.RespNFTTransferDetailsByAddress, 0)
	offset := (req.Page - 1) * req.PageSize
	offsetStr := strconv.Itoa(int(offset))

	orm := orm.NewOrm()
	condition := " WHERE token_address = '" + req.ContractAddress + "' and (`from`='" + req.Address + "' OR `to` = '" + req.Address + "') GROUP BY token_id "
	sqlStr := "SELECT max(tx_time), token_id, `from`, `to`, count(0) as transfer_count from tb_contract_transaction " + condition + " LIMIT " + req.Length + " OFFSET " + offsetStr
	_, err := orm.Raw(sqlStr).QueryRows(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// NFTTxDetail NFT溯源-交易详情
func NFTTxDetail(req apiModels.ReqNFTTxDetail) (*apiModels.RespNFTTxDetail, int, error) {
	var res *apiModels.RespNFTTxDetail
	orm := orm.NewOrm()

	var count int
	sqlCount := "SELECT count(0) as transfer_count from tb_contract_transaction where tx_hash = '" + req.TxHash + "' and token_address = '" + req.ContractAddress +"'"
	orm.Raw(sqlCount).QueryRow(&count)

	sqlStr := "SELECT token_id, `from`, `to`, max(tx_time) as tx_time, topics1 as method from tb_contract_transaction c left join tb_event e on c.tx_hash = e.tx_hash where c.tx_hash = '" + req.TxHash + "' and token_address = '" + req.ContractAddress +"'"
	err := orm.Raw(sqlStr).QueryRow(&res)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}


// NFTTransferDetailByTokenId NFT溯源-交易详情-根据token id获取流转详情列表
func NFTTransferDetailByTokenId(req apiModels.ReqNFTTransferDetailsByTokenId) ([]*apiModels.RespNFTTransferDetailsByTokenId, error) {
	res := make([]*apiModels.RespNFTTransferDetailsByTokenId, 0)
	orm := orm.NewOrm()

	offset := (req.Page - 1) * req.PageSize
	offsetStr := strconv.Itoa(int(offset))

	//token_id 转换
	tokenId, ok := new(big.Int).SetString(req.TokenID, 10)
	if !ok {
		return nil, errors.New("convert token_id error.")
	}
	req.TokenID = common.BigToHash(tokenId).String()

	sqlStr := "SELECT c.`from`, c.`to`, c.tx_time, t.input_data as method from tb_contract_transaction c join tb_transaction t on c.tx_hash = t.tx_hash where c.token_id = '" + req.TokenID + "' and token_address= '" + req.ContractAddress + "'" +
		" LIMIT " + req.Length + " OFFSET " + offsetStr
	_, err := orm.Raw(sqlStr).QueryRows(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// NFTDetail NFT溯源-NFT详情
func NFTDetail(req apiModels.ReqNFTDetail) (*apiModels.RespNFTDetail, error) {
	var res *apiModels.RespNFTDetail
	orm := orm.NewOrm()

	var txCount int
	countSql := "SELECT count(0) as transfer_count from tb_contract_transaction where token_address = '" + req.ContractAddress + "' and token_id = '" + req.TokenID + "'"
	err := orm.Raw(countSql).QueryRow(&txCount)

	var histHolderCount int
	histHoldersSql := "SELECT count(distinct `to`) as history_holder_count from tb_contract_transaction where token_address = '" + req.ContractAddress + "' and token_id = '" + req.TokenID +"'"
	err = orm.Raw(histHoldersSql).QueryRow(&histHolderCount)

	var mintTime time.Time
	mintTimeSql := "SELECT tx_time from tb_contract_transaction where token_address = '" + req.ContractAddress + "' and token_id = '" + req.TokenID +"' and `from` = '0x0000000000000000000000000000000000000000'"
	err = orm.Raw(mintTimeSql).QueryRow(&mintTime)
	if err != nil {
		return nil, err
	}

	res.TransferCount = txCount
	res.HistoryHolderCount = histHolderCount

	// 处理mint_time
	res.MintTime = utils.Float64String(time.Since(mintTime).Hours() / 24)
	return res, nil
}

// LongestHold NFT溯源-地址详情-最长持有的tokenId和持有时间
// todo 最好改为查nebula
func LongestHold(contractAddress, accountAddress string) ([]*TbContractTransaction, error) {
	res := make([]*TbContractTransaction, 0)
	orm := orm.NewOrm()
	sqlStr := "SELECT token_id, tx_time,`from`,`to` FROM tb_contract_transaction WHERE token_address = '" + contractAddress + "' and (`from` = '" + accountAddress + "' or `to` = '" + accountAddress + "') order by tx_time asc"
	_, err := orm.Raw(sqlStr).QueryRows(&res)
	return res, err
}

