package service

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/server/txs-analysis/constant"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/models/dbModels"
	"github.com/server/txs-analysis/models/nebulaModels"
)

type AddressService struct {
}

// AddressTxList 可查的token下拉列表
func (this *AddressService) AddressTokenList(req apiModels.ReqContractList) ([]*apiModels.RespContractList, error) {
	contracts, err := dbModels.QueryContractList(req.TxType, req.Name)
	if err != nil && err != orm.ErrNoRows {
		beego.Error("dbModels.QueryContractList error.", err)
	}
	return contracts, nil
}

// AddressDetail 地址分析-地址详情
func (this *AddressService) AddressDetail(address, contractAddress string) (*apiModels.RespAddressDetail, error) {
	var accountInfo *apiModels.RespContractAddressInfo
	var err error
	// CMP
	if len(contractAddress) == 0 {
		accountInfo, err = dbModels.GetAddressInfo(address)
		if err == orm.ErrNoRows {
			beego.Error("dbModels.GetAddressInfo error.", err, " address:", address)
			return nil, errors.New(constant.ErrAddress)
		}
		if err != nil {
			beego.Error("dbModels.GetAddressInfo error.", err, " address:", address)
		}

	} else {
		accountInfo, err = dbModels.GetContractAddressInfo(address, contractAddress)
		if err != nil && err != orm.ErrNoRows {
			beego.Error("dbModels.GetContractAddressInfo error.", err, " address:", address, "contract_address:", contractAddress)
		}
	}

	var res apiModels.RespAddressDetail
	res.Address = address

	if accountInfo != nil {
		res.Type = accountInfo.AccountType
		res.Balance = accountInfo.Balance
	}

	// 地址的交易相关信息
	addressTxInfo, err := dbModels.GetAddressTxInfo(address, contractAddress)
	if err != nil && err != orm.ErrNoRows {
		beego.Error("dbModels.GetAddressTxInfo error.", err)
		return nil, err
	}

	//组装数据
	if addressTxInfo != nil {
		res.OutAddressCount = addressTxInfo.OutAddressCount
		res.InAddressCount = addressTxInfo.InAddressCount
		res.FirstTxTime = addressTxInfo.FirstTxTime
		res.TxCount = addressTxInfo.TxCount
		res.MaxTxAmount = addressTxInfo.MaxTxAmount
		res.ReceiveAmountTotal = addressTxInfo.ReceiveAmountTotal
		res.SendAmountTotal = addressTxInfo.SendAmountTotal
	}
	return &res, nil
}

// AddressTxDetail 地址分析-交易详情
func (this *AddressService) AddressTxDetail(req apiModels.ReqAddressTxDetail) (*apiModels.RespTxDetailInfo, error) {
	txDetailInfo, err := dbModels.GetAddressTxDetailInfo(req)
	if err != nil {
		beego.Error("dbModels.GetAddressTxDetailInfo error.", err)
	}
	if txDetailInfo != nil {
		fromType, err := dbModels.GetAddressType(req.From)
		if err != nil {
			beego.Error("from address dbModels.GetAddressType error.", err)
		}

		toType, err := dbModels.GetAddressType(req.To)
		if err != nil {
			beego.Error("to address dbModels.GetAddressType error.", err)
		}

		txDetailInfo.FromType = fromType
		txDetailInfo.ToType = toType
		txDetailInfo.From = req.From
		txDetailInfo.To = req.To
	}
	return txDetailInfo, nil
}

// AddressTxList 地址分析-交易列表
func (this *AddressService) AddressTxList(req apiModels.ReqAddressTxList) ([]*apiModels.RespAddressTxList, error) {
	var res []*apiModels.RespAddressTxList
	var err error
	if len(req.ContractAddress) == 0 {
		res, err = dbModels.GetAddressTxDetailList(req)
		if err != nil {
			beego.Error("dbModels.GetAddressTxDetailList error.", err)
		}
	} else {
		res, err = dbModels.GetAddressTokenTxDetailList(req)
		if err != nil {
			beego.Error("dbModels.GetAddressTokenTxDetailList error.", err)
		}
	}
	return res, nil
}

// AddressTxAnalysis 地址分析-地址交易图
func (this *AddressService) AddressTxAnalysis(req apiModels.ReqAddressTxGraph) (*nebulaModels.RespGraph, error) {
	res, err := nebulaModels.AddressTxAnalysis(req)
	if err != nil {
		beego.Error("nebulaModels.AddressTxAnalysis error.", err)
	}
	return res, nil
}
