package dbModels

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type ResSyncTransaction struct {
	BlockId         string    `json:"block_id"`
	TxHash          string    `json:"tx_hash"`
	TxTime          time.Time `json:"tx_time"`
	TxTimeStr       string    `json:"tx_time_str"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	Amount          string    `json:"amount"`
	Value           string    `json:"value"`
	TokenType       string    `json:"token_type"`
	ContractAddress string    `json:"contract_address"`
	TokenAddress    string    `json:"token_address"`
	TokenId         string    `json:"token_id"`
	Caller          string    `json:"caller"`
}

//获取tb_transaction的CMP交易数据
func GetSyncTxData(start, end string) (Res []*ResSyncTransaction, err error) {
	ormer := orm.NewOrm()
	list := make([]*ResSyncTransaction, 0)
	sql := "select * from tb_transaction where `value` !='0x0' and tx_status = 1 and block_id >=" + start + " and block_id <=" + end + " order by tx_time asc"
	_, err = ormer.Raw(sql).QueryRows(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

//获取tb_internal_transaction表交易金额>0的交易
func GetSyncInternalTxData(start, end string) (Res []*ResSyncTransaction, err error) {
	ormer := orm.NewOrm()
	list := make([]*ResSyncTransaction, 0)
	sql := "select A.*,B.from as caller from tb_internal_transaction A left join tb_transaction B on A.tx_hash = B.tx_hash where LENGTH(A.value)>0 and A.value !='0x0' and A.tx_status = 1 and A.block_id >=" + start + " and A.block_id <=" + end + " order by tx_time asc"
	_, err = ormer.Raw(sql).QueryRows(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

//获取token交易
func GetSyncTokenTxsData(start, end string) (Res []*ResSyncTransaction, err error) {
	ormer := orm.NewOrm()
	list := make([]*ResSyncTransaction, 0)
	sql := "select A.*,B.from as caller from tb_contract_transaction A left join tb_transaction B on A.tx_hash=B.tx_hash where A.token_type = 1 and A.tx_status = 1 and A.block_id >=" + start + " and A.block_id <=" + end + " order by tx_time asc"
	_, err = ormer.Raw(sql).QueryRows(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

//获取nft交易
func GetSyncNFTTxsData(start, end string) (Res []*ResSyncTransaction, err error) {
	ormer := orm.NewOrm()
	list := make([]*ResSyncTransaction, 0)
	sql := "select A.*,B.from as caller from tb_contract_transaction A left join tb_transaction B on A.tx_hash=B.tx_hash where (A.token_type = 2 or A.token_type = 4) and A.tx_status = 1 and A.block_id >=" + start + " and A.block_id <=" + end + " order by tx_time asc"
	_, err = ormer.Raw(sql).QueryRows(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}
