package dbModels

import (
	"github.com/server/txs-analysis/models/apiModels"
	"time"

	"github.com/astaxie/beego/orm"
)

type TbTransaction struct {
	Id                int64     `orm:"column(id);pk"`
	BlockId           int64     `orm:"column(block_id);null" description:"区块号"`
	BlockHash         string    `orm:"column(block_hash);size(255);null" description:"区块hash"`
	TxHash            string    `orm:"column(tx_hash);size(255);null" description:"交易hash"`
	TxType            int       `orm:"column(tx_type);null" description:"交易类型 1：普通交易 2-合约交易"`
	TxTime            time.Time `orm:"column(tx_time);type(datetime);null" description:"交易时间"`
	TxIndex           int64     `orm:"column(tx_index);null" description:"交易索引"`
	ContractAddress   string    `orm:"column(contract_address);size(255);null" description:"合约地址"`
	From              string    `orm:"column(from);size(255);null" description:"发起地址"`
	To                string    `orm:"column(to);size(255);null" description:"到大地址"`
	Value             string    `orm:"column(value);size(255);null" description:"交易值"`
	Amount            float64   `orm:"column(amount);null;digits(64);decimals(18)" description:"交易金额"`
	TxFee             string    `orm:"column(tx_fee);size(255);null" description:"交易费"`
	TxStatus          int       `orm:"column(tx_status);null" description:"交易状态"`
	ErrorInfo         string    `orm:"column(error_info);null" description:"错误信息"`
	CumulativeGasUsed int64     `orm:"column(cumulative_gas_used);null" description:"CumulativeGasUsed"`
	GasUsed           int64     `orm:"column(gas_used);null" description:"使用gas"`
	GasLimit          int64     `orm:"column(gas_limit);null" description:"限制gas"`
	GasPrice          int64     `orm:"column(gas_price);null" description:"gas价格"`
	Nonce             int64     `orm:"column(nonce);null" description:"nonce号"`
	InputData         string    `orm:"column(input_data);null" description:"输入数据"`
	Status            int       `orm:"column(status);null" description:"状态,1-正常 2-冻结"`
	IsDeleted         int8      `orm:"column(is_deleted);null" description:"删除状态 0-正常 1-删除"`
	SyncTime          time.Time `orm:"column(sync_time);type(timestamp);null;auto_now_add" description:"同步时间"`
	CreateTime        time.Time `orm:"column(create_time);type(timestamp);null;auto_now_add" description:"创建时间"`
	UpdateTime        time.Time `orm:"column(update_time);type(timestamp);null;auto_now_add" description:"更新时间"`
}

func (t *TbTransaction) TableName() string {
	return "tb_transaction"
}

func init() {
	orm.RegisterModel(new(TbTransaction))
}

// 地址分析-地址详情
func GetAddressTxInfo(address, contractAddr string) (Res *apiModels.RespAddressDetail, err error) {
	ormer := orm.NewOrm()

	var table string
	var condition string
	if len(contractAddr) == 0 {
		table = "tb_transaction"
		condition = ""
	} else {
		table = "tb_contract_transaction"
		condition = " and token_address = '" + contractAddr + "'"
	}

	sqlStr := " select A.address, sum(tx_count) as tx_count, sum(in_addr_count) as in_address_count, sum(out_addr_count) as out_address_count, min(min_tx_time) as first_tx_time, max(max_amount) as max_tx_amount, sum(in_amount) as receive_amount_total, sum(out_amount) as send_amount_total " +
		" from ( " +
		" select `from` as address, COUNT(1) as tx_count, COUNT(DISTINCT `from`) as in_addr_count, 0 as out_addr_count, sum(amount) as in_amount, 0.0 as out_amount, min(tx_time) as min_tx_time, max(amount) as max_amount from " + table + " where`to` = '" + address + "'" + condition + " GROUP BY `from` " +
		" union " +
		" select `to` as address, COUNT(1) as tx_count, 0 as in_addr_count, COUNT(DISTINCT `to`) as out_addr_count, 0.0 as in_amount, sum(amount) as out_amount, min(tx_time) as min_tx_time, max(amount) as max_amount from " + table + " t where `from` = '" + address + "'" + condition + " GROUP BY `to` " +
		" )A group by A.address "

	err = ormer.Raw(sqlStr).QueryRow(&Res)
	return
}

//func GetAddressTxInfo(address string) (Res *apiModels.RespAddressDetail, err error) {
//	ormer := orm.NewOrm()
//	sqlStr := " select A.address, sum(tx_count) as tx_count, sum(in_addr_count) as in_address_count, sum(out_addr_count) as out_address_count, min(min_tx_time) as first_tx_time, max(max_amount) as max_tx_amount, sum(in_amount) as receive_amount_total, sum(out_amount) as send_amount_total " +
//		" from ( " +
//		" select `from` as address, COUNT(1) as tx_count, COUNT(DISTINCT `from`) as in_addr_count, 0 as out_addr_count, sum(amount) as in_amount, 0.0 as out_amount, min(tx_time) as min_tx_time, max(amount) as max_amount from tb_transaction where`to` = '" + address + "' GROUP BY `from` " +
//		" union " +
//		" select `to` as address, COUNT(1) as tx_count, 0 as in_addr_count, COUNT(DISTINCT `to`) as out_addr_count, 0.0 as in_amount, sum(amount) as out_amount, min(tx_time) as min_tx_time, max(amount) as max_amount from tb_transaction t where `from` = '" + address + "' GROUP BY `to` " +
//		" )A group by A.address "
//
//	err = ormer.Raw(sqlStr).QueryRow(&Res)
//	return
//}

// 地址分析-地址交易图数据
func GetAddressTxList(address string) (Res *apiModels.RespAddressDetail, err error) {
	ormer := orm.NewOrm()
	sqlStr := " select A.address, sum(in_tx_count) as in_tx_count, sum(out_tx_count) as out_tx_count, sum(in_addr_count) as in_address_count, sum(out_addr_count) as out_address_count, min(min_tx_time) as first_tx_time, max(max_amount) as max_tx_amount, sum(in_amount) as receive_amount_total, sum(out_amount) as send_amount_total " +
		" from ( " +
		" select `from` as address,count(1) as in_tx_count, 0 as out_tx_count, COUNT((DISTINCT `from`) as in_addr_count, 0 as out_addr_count, sum(amount) as in_amount, 0.0 as out_amount, min(tx_time) as min_tx_time, max(amount) as max_amount from tb_transaction where`to` = '" + address + "' GROUP BY `from` " +
		" UNION " +
		" select to as address, 0 as in_tx_count, count(1) as out_tx_count, 0 as in_addr_count, COUNT(DISTINCT `to`) as out_addr_count, 0.0 as in_amount, sum(amount) as out_amount from tb_transaction t where `from` = '" + address + "' GROUP BY `to` " +
		" )A group by A.address "

	//Res.TxCount = Res.InTxCount + Res.OutTxCount
	_, err = ormer.Raw(sqlStr).QueryRows(&Res)
	return
}

// 地址分析-交易详情-交易信息
func GetAddressTxDetailInfo(req apiModels.ReqAddressTxDetail) (Res *apiModels.RespTxDetailInfo, err error) {
	ormer := orm.NewOrm()
	condition := ""
	table := "tb_transaction"
	if len(req.ContractAddress) > 0 {
		table = "tb_contract_transaction"
		condition = " and token_address = '" + req.ContractAddress + "'"
	}
	sqlStr := "SELECT count(distinct tx_hash) as tx_count, sum(amount) as tx_amount, min(tx_time) as first_tx_time, max(tx_time) as latest_tx_time from " + table + " WHERE `from`='" + req.From + "' and `to`= '" + req.To + "'" + condition
	err = ormer.Raw(sqlStr).QueryRow(&Res)
	return
}

// 地址分析-交易详情-交易列表
func GetAddressTxDetailList(fromAddress, toAddress string) (Res []*apiModels.RespAddressTxList, err error) {
	ormer := orm.NewOrm()
	sqlStr := "SELECT tx_time, tx_hash, amount from tb_transaction WHERE `from`='" + fromAddress + "' and `to`= '" + toAddress + "'"
	_, err = ormer.Raw(sqlStr).QueryRows(&Res)
	return
}

// 交易图谱-交易详情
func GetTxInfo(txHash string) (Res *apiModels.RespTxDetail, err error) {
	ormer := orm.NewOrm()
	sqlStr := "SELECT tx_hash, amount, tx_fee, tx_time,`from`,`to` FROM tb_transaction WHERE tx_hash='" + txHash + "';"
	err = ormer.Raw(sqlStr).QueryRow(&Res)
	return
}

// 交易图谱-地址详情
//地址类型、余额、流出金额、流出金额 流入/流出地址数 流入/流出笔数
func GetTxAddressDetail(address string) (Res *apiModels.RespTxAddressDetail, err error) {
	ormer := orm.NewOrm()
	sqlStr := " select A.address, sum(in_tx_count) as in_tx_count, sum(out_tx_count) as out_tx_count, sum(in_addr_count) as in_address_count, sum(out_addr_count) as out_address_count, sum(in_amount) as receive_amount_total, sum(out_amount) as send_amount_total " +
		" from ( " +
		" select `from` as address, count(1) as in_tx_count, 0 as out_tx_count, COUNT(DISTINCT `from`) as in_addr_count, 0 as out_addr_count, sum(amount) as in_amount, 0.0 as out_amount from tb_transaction where`to` = '" + address + "' GROUP BY `from` " +
		" UNION " +
		" select `to` as address, 0 as in_tx_count, count(1) as out_tx_count, 0 as in_addr_count, COUNT(DISTINCT `to`) as out_addr_count, 0.0 as in_amount, sum(amount) as out_amount from tb_transaction t where `from` = '" + address + "' GROUP BY `to` " +
		" )A group by A.address "

	ormer.Raw(sqlStr).QueryRow(&Res)
	return
}

// 交易图谱-交易图数据
func TxGraphData(req apiModels.ReqAddressTxGraph) (Res []*apiModels.RespAddressDetail, err error) {
	ormer := orm.NewOrm()

	limit := req.Count
	address := req.Address
	reqType := req.Direction
	sqlStr := ""
	condition := ""
	table := "tb_transaction"
	if len(req.ContractAddress) > 0 {
		condition = " and contract_address = '" + req.ContractAddress + "'"
		table = "tb_contract_transaction"
	}

	if reqType == "all" {
		sqlStr = "select A.address, sum(tx_count) as tx_count, sum(amount) as tx_amount from " +
			" (select `from` as address, count(1) as tx_count, sum(amount) as amount from " + table + " where`to` = '" + address + "'" + condition + " GROUP BY `from` " +
			" union " +
			" select `to` as address, count(1) as tx_count, sum(amount) as amount from " + table + " t where `from` = '" + address + "'" + condition + " GROUP BY `to` " +
			" )A group by A.address"
	} else {
		addrType := ""
		if reqType == "out" {
			addrType = "`from`"
		} else {
			addrType = "`to`"
		}
		sqlStr = " select count(1) as tx_count, sum(amount) as tx_amount from " + table + " where " + addrType + " = '" + address + "'" + condition + "' GROUP BY " + addrType + " limit " + limit
	}
	ormer.Raw(sqlStr).QueryRows(&Res)
	return
}
