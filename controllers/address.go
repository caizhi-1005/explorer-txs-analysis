package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/server/txs-analysis/constant"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/service"
)

type AddressController struct {
	BaseController
	addressService service.AddressService
}

// AddressTokenList 可查的token下拉列表
func (this *AddressController) AddressTokenList() {
	this.IsPost()
	Req := apiModels.ReqContractList{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}

	res, err := this.addressService.AddressTokenList(Req)
	if err != nil {
		beego.Error(constant.ErrSystem, err)
		this.ResponseInfo(500, constant.ErrSystem, nil)
	}
	this.ResponseInfo(200, nil, res)
}

// AddressDetail 地址分析-地址详情
func (this *AddressController) AddressDetail() {
	this.IsPost()
	Req := apiModels.ReqAddressDetail{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	if len(Req.Field) <= 0 || len(Req.Value) <= 0 {
		beego.Error(constant.ErrParam)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	res, err := this.addressService.AddressDetail(Req.Value, Req.ContractAddress)
	if err != nil {
		beego.Error(constant.ErrSystem, err)
		this.ResponseInfo(500, constant.ErrSystem, err.Error())
	}
	this.ResponseInfo(200, nil, res)
}

// AddressTxDetail 地址分析-交易详情
func (this *AddressController) AddressTxDetail() {
	this.IsPost()
	Req := apiModels.ReqAddressTxDetail{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	if len(Req.From) <= 0 || len(Req.To) <= 0 {
		beego.Error(constant.ErrParam)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	res, err := this.addressService.AddressTxDetail(Req)
	if err != nil {
		beego.Error(constant.ErrSystem, err)
		this.ResponseInfo(500, constant.ErrSystem, nil)
	}
	this.ResponseInfo(200, nil, res)
}

// AddressTxList 地址分析-交易列表
func (this *AddressController) AddressTxList() {
	this.IsPost()
	Req := apiModels.ReqAddressTxList{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}

	page, pageSize := this.Pagination(Req.Start, Req.Length)
	Req.Page = page
	Req.PageSize = pageSize

	res, err := this.addressService.AddressTxList(Req)
	if err != nil {
		beego.Error(constant.ErrSystem, err)
		this.ResponseInfo(500, constant.ErrSystem, nil)
	}
	this.ResponseInfo(200, nil, res)
}

// AddressTxAnalysis 地址分析-地址交易图
func (this *AddressController) AddressTxAnalysis() {
	this.IsPost()
	//Req := apiModels.ReqTxAnalysis{}
	Req := apiModels.ReqAddressTxGraph{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	if len(Req.Address) == 0 {
		beego.Error(constant.ErrParam)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	res, err := this.addressService.AddressTxAnalysis(Req)
	if err != nil {
		beego.Error(constant.ErrSystem, err)
		this.ResponseInfo(500, constant.ErrSystem, nil)
	}
	this.ResponseInfo(200, nil, res)
}
