package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/models/dbModels"
)

type AddressService struct {
}

// AddressDetail 地址分析-地址详情
func (this *AddressService) AddressDetail(address string) (*apiModels.RespAddressDetail, error) {
	// 地址类型
	accountInfo, err := dbModels.GetAddressInfo(address)
	if err != nil && err != orm.ErrNoRows {
		beego.Error("dbModels.GetAddressType error.", err)
	}

	// 地址的交易相关信息
	addressTxInfo, err := dbModels.GetAddressTxInfo(address)
	if err != nil && err != orm.ErrNoRows {
		beego.Error("dbModels.GetAddressTxInfo error.", err)
		return nil, err
	}

	//组装数据
	if addressTxInfo != nil {
		res := &apiModels.RespAddressDetail{
			Address : address,
			Type : accountInfo.AccountType,
			Balance : accountInfo.Balance,
			OutAddressCount : addressTxInfo.OutAddressCount,
			InAddressCount : addressTxInfo.InAddressCount,
			FirstTxTime : addressTxInfo.FirstTxTime,
			TxCount : addressTxInfo.TxCount,
			MaxTxAmount : addressTxInfo.MaxTxAmount,
			ReceiveAmountTotal : addressTxInfo.ReceiveAmountTotal,
			SendAmountTotal : addressTxInfo.SendAmountTotal,
		}
		return res, nil
	}
	return nil, nil
}

// AddressTxAnalysis 地址分析-地址交易图
func (this *AddressService) AddressTxAnalysis(req apiModels.ReqTxAnalysis) ([]*apiModels.RespAddressDetail, error) {
	//accountInfo, err := dbModels.GetAddressTxList(address)

	reqGraph := apiModels.ReqAddressTxGraph{}
	reqGraph.Address = req.Address
	reqGraph.Type = "0"
	res, err := dbModels.TxGraphData(reqGraph)
	if err != nil {
		beego.Error("dbModels.TxGraphData error.", err)
	}
	return res, nil
}

// AddressTxDetail 地址分析-交易详情
func (this *AddressService) AddressTxDetail(fromAddress, toAddress string) (*apiModels.RespTxDetailInfo, error) {
	txDetailInfo, err := dbModels.GetAddressTxDetailInfo(fromAddress, toAddress)
	if err != nil {
		beego.Error("dbModels.GetAddressTxDetailInfo error.", err)
	}
	return txDetailInfo, nil
}

// AddressTxList 地址分析-交易列表
func (this *AddressService) AddressTxList(fromAddress, toAddress string) ([]*dbModels.TbTransaction, error) {
	tbTransactionList, err := dbModels.GetAddressTxDetailList(fromAddress, toAddress)
	if err != nil {
		beego.Error("dbModels.GetAddressTxDetailList error.", err)
	}
	return tbTransactionList, nil
}
