package dbModels

import (
	"github.com/server/txs-analysis/models/apiModels"
	"time"

	"github.com/astaxie/beego/orm"
)

type TbContractAccount struct {
	Id              int64     `orm:"column(id);pk"`
	BlockId         int64     `orm:"column(block_id);null" description:"区块号"`
	BlockHash       string    `orm:"column(block_hash);size(255);null" description:"blockHash"`
	TxHash          string    `orm:"column(tx_hash);size(255);null" description:"交易hash"`
	ContractAddress string    `orm:"column(contract_address);size(255);null" description:"合约地址"`
	ContractType    int       `orm:"column(contract_type);null" description:"账户类型：1-erc20 2-erc721"`
	AccountAddress  string    `orm:"column(account_address);size(255);null" description:"账户地址"`
	AccountType     int       `orm:"column(account_type);null" description:"账户类型：1-account 2-contract"`
	Balance         float64   `orm:"column(balance);null;digits(64);decimals(18)" description:"余额"`
	TokenId         string    `orm:"column(token_id);size(255);null" description:"tokenId"`
	Status          int       `orm:"column(status);null" description:"状态,1-正常 2-冻结"`
	IsDeleted       int8      `orm:"column(is_deleted);null" description:"删除状态 0-正常 1-删除"`
	SyncTime        time.Time `orm:"column(sync_time);type(timestamp);null;auto_now_add" description:"同步时间"`
	CreateTime      time.Time `orm:"column(create_time);type(timestamp);null;auto_now_add" description:"创建时间"`
	UpdateTime      time.Time `orm:"column(update_time);type(timestamp);null;auto_now_add" description:"更新时间"`
}

func (t *TbContractAccount) TableName() string {
	return "tb_contract_account"
}

func init() {
	orm.RegisterModel(new(TbContractAccount))
}

// GetContractAddressInfo 根据合约地址和账户地址获取账户信息
func GetContractAddressInfo(address, contractAddress string) (*apiModels.RespContractAddressInfo, error) {
	ormer := orm.NewOrm()
	var accountInfo *apiModels.RespContractAddressInfo
	sql := "select a.account_type, a.balance, c.symbol from tb_contract_account a left join tb_contract c on a.contract_address = c.contract_address where c.name != '' and c.symbol != '' and a.account_address ='" + address + "' and a.contract_address = '" + contractAddress + "'"
	err := ormer.Raw(sql).QueryRow(&accountInfo)
	if err != nil {
		return nil, err
	}
	return accountInfo, nil
}


// HoldTokenIdAndAddressType NFT溯源-地址详情-获取地址类型和持有的token Id
func HoldTokenIdCount(contractAddress, accountAddress string) (int, error) {
	var count int
	orm := orm.NewOrm()
	countSql := "SELECT count(0) FROM tb_contract_account WHERE contract_address = '" + contractAddress + "' and account_address = '" + accountAddress +"'"
	err := orm.Raw(countSql).QueryRow(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}


// HoldTokenIdAndAddressType NFT溯源-地址详情-获取地址类型和持有的token Id
func HoldTokenIdAndAddressType(contractAddress, accountAddress string) ([]apiModels.HoldTokenId, error) {
	list := make([]apiModels.HoldTokenId, 0)
	orm := orm.NewOrm()
	sqlStr := "SELECT account_type, distinct(token_id) FROM tb_contract_account WHERE contract_address = '" + contractAddress + "' and account_address = '" + accountAddress +"'"
	_, err := orm.Raw(sqlStr).QueryRows(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// HoldTokenIdHistory NFT溯源-地址详情-历史持有的tokenId列表
func HoldTokenIdHistory(contractAddress, accountAddress string) ([]apiModels.HoldTokenIdHistory, error) {
	list := make([]apiModels.HoldTokenIdHistory, 0)
	orm := orm.NewOrm()
	sqlStr := "SELECT distinct(token_id) FROM tb_contract_transaction WHERE token_address = '" + contractAddress + "' and `to` = '" + accountAddress +"'"
	_, err := orm.Raw(sqlStr).QueryRows(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}


