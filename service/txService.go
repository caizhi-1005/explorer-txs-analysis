package service

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/server/txs-analysis/constant"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/models/dbModels"
	"github.com/server/txs-analysis/utils"
)

type TxService struct {
}

// TxService 交易图谱-交易详情
func (this *TxService) TxDetail(req apiModels.ReqTxDetail) (*apiModels.RespTxDetail, error) {
	tx, err := dbModels.GetTxInfo(req)
	if err == orm.ErrNoRows {
		beego.Error("dbModels.GetTxInfo error.", err, " tx_hash:", req.Value)
		return nil, errors.New(constant.ErrTxHash)
	}
	if err != nil {
		beego.Error("dbModels.GetTxInfo error.", err, " tx_hash:", req.Value)
	}
	if tx != nil {
		//处理tx_fee小数位
		tx.TxFee = utils.FeeFormatToDecimalAmount(tx.TxFee, constant.BASE_TOKEN_DECIMAL)
	}
	return tx, nil
}


// TxAddressDetail 交易图谱-地址详情
func (this *TxService) TxAddressDetail(req apiModels.ReqTxDetail) (*apiModels.RespTxAddressDetail, error) {
	detail, err := dbModels.GetTxAddressDetail(req)
	if err == orm.ErrNoRows {
		beego.Error("dbModels.GetTxAddressDetail error.", err, " address:", req.Value)
		return nil, errors.New(constant.ErrAddress)
	}
	if err != nil {
		beego.Error("dbModels.GetTxAddressDetail error.", err)
	}
	return detail, nil
}

// TxGraphData 交易图谱-交易图
func (this *TxService) TxGraphData(req apiModels.ReqAddressTxGraph) ([]*apiModels.RespAddressDetail, error) {
	detail, err := dbModels.TxGraphData(req)
	if err != nil {
		beego.Error("dbModels.TxGraphData error.", err)
	}
	return detail, nil
}
