package dbModels

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type TbAccountInfo struct {
	Id             int64     `orm:"column(id);pk"`
	BlockId        int64     `orm:"column(block_id);null" description:"区块号"`
	BlockHash      string    `orm:"column(block_hash);size(255);null" description:"区块哈希"`
	TxHash         string    `orm:"column(tx_hash);size(255);null" description:"交易哈希"`
	AccountAddress string    `orm:"column(account_address);size(255);null" description:"账户地址"`
	AccountType    int       `orm:"column(account_type);null" description:"账户类型：1-account 2-contract"`
	Balance        float64   `orm:"column(balance);null;digits(64);decimals(18)" description:"余额"`
	SyncTime       time.Time `orm:"column(sync_time);type(timestamp);null;auto_now_add" description:"同步时间"`
}

func (t *TbAccountInfo) TableName() string {
	return "tb_account_info"
}

func init() {
	orm.RegisterModel(new(TbAccountInfo))
}

type ResWithdrawAccountList struct {
	AccountAddress string `json:"account_address"`
	AccountType    int    `json:"account_type"`
}

type ResTxListByAddress struct {
	TxList        interface{} `json:"tx_list"`
	InTotalValue  float64     `json:"in_total_value"`
	OutTotalValue float64     `json:"out_total_value"`
}

func GetWithdrawAccountList() (Res []*ResWithdrawAccountList, err error) {
	ormer := orm.NewOrm()
	list := make([]*ResWithdrawAccountList, 0)
	//sq1 := "select A.from as account_address,B.account_type,B.balance from tb_transaction A left join tb_account_info B on A.from=B.account_address where A.tx_status = 1 and A.input_data='0x86d1a69f' or A.`to`= '0x871fcb6b836db1b5d6ee64901fb17245cd403e6d'"
	sql := "SELECT DISTINCT(`from`) as account_address FROM tb_transaction WHERE tx_status = 1 AND input_data = '0x86d1a69f' OR `to` = '0x871fcb6b836db1b5d6ee64901fb17245cd403e6d'"
	_, err = ormer.Raw(sql).QueryRows(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func AccountInfoList(filters ...interface{}) ([]*TbAccountInfo, error) {
	list := make([]*TbAccountInfo, 0)
	query := orm.NewOrm().QueryTable(new(TbAccountInfo).TableName())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	_, err := query.All(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func GetSyncAddressData(start, end string) ([]*TbAccountInfo, error) {
	ormer := orm.NewOrm()
	list := make([]*TbAccountInfo, 0)
	sql := "select account_address,account_type from tb_account_info where block_id >=" + start + " and block_id <=" + end
	_, err := ormer.Raw(sql).QueryRows(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}


func GetAddressInfo(address string) (*TbAccountInfo, error) {
	ormer := orm.NewOrm()
	var accountInfo *TbAccountInfo
	//accountInfo := new(TbAccountInfo)
	sql := "select account_type, balance from tb_account_info where account_address ='" + address + "'"
	err := ormer.Raw(sql).QueryRow(&accountInfo)
	if err != nil {
		return nil, err
	}
	return accountInfo, nil
}
