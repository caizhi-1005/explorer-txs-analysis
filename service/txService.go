package service

import (
	"github.com/astaxie/beego"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/models/dbModels"
)

type TxService struct {
}


// TxService 交易图谱-交易详情
//func (this *TxService) TxDetail(address string) (*apiModels.RespTxDetail, error) {
//	var res *apiModels.RespTxDetail
//	tx, err := dbModels.GetTxInfo(address)
//	if err != nil {
//		beego.Error("dbModels.GetTxInfo error.", err)
//	}
//	//获取地址类型
//	fromInfo, err := dbModels.GetAddressInfo(tx.From)
//	if err != nil {
//		beego.Error("dbModels.GetAddressInfo error.", err)
//	}
//	toInfo, err := dbModels.GetAddressInfo(tx.To)
//	if err != nil {
//		beego.Error("dbModels.GetAddressInfo error.", err)
//	}
//
//	res.FromType = fromInfo.AccountType
//	res.ToType = toInfo.AccountType
//	return tx, nil
//}

// TxService 交易图谱-交易详情
func (this *TxService) TxDetail(address string) (*apiModels.RespTxDetail, error) {
	tx, err := dbModels.GetTxInfo(address)
	if err != nil {
		beego.Error("dbModels.GetTxInfo error.", err)
	}
	return tx, nil
}


// TxAddressDetail 交易图谱-地址详情
func (this *TxService) TxAddressDetail(address string) (*apiModels.RespTxAddressDetail, error) {
	detail, err := dbModels.GetTxAddressDetail(address)
	if err != nil {
		beego.Error("dbModels.GetTxAddressDetail error.", err)
	}
	return detail, nil
}

// TxGraphData 交易图谱-交易图
func (this *TxService) TxGraphData(req apiModels.ReqAddressTxGraph) ([]*apiModels.RespAddressDetail, error) {
	//limit := req.Count
	//address := req.Address

	detail, err := dbModels.TxGraphData(req)
	if err != nil {
		beego.Error("dbModels.GetTxAddressDetail error.", err)
	}
	return detail, nil
}
