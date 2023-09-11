package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/server/txs-analysis/constant"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/service"
)

type TxController struct {
	BaseController
	txService service.TxService
}

// TxDetail 交易图谱-交易详情
func (this *TxController) TxDetail() {
	this.IsPost()
	Req := apiModels.ReqTxDetail{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	if len(Req.Field) <= 0 || len(Req.Value) <= 0 {
		beego.Error(constant.ErrParam)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	res, err := this.txService.TxDetail(Req)
	if err != nil {
		beego.Error(constant.ErrSystem, err)
		this.ResponseInfo(500, constant.ErrSystem, nil)
	}
	this.ResponseInfo(200, nil, res)
}


// TxAddressDetail 交易图谱-地址详情
func (this *TxController) TxAddressDetail() {
	this.IsPost()
	Req := apiModels.ReqTxDetail{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	if len(Req.Field) <= 0 || len(Req.Value) <= 0 {
		beego.Error(constant.ErrParam)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	res, err := this.txService.TxAddressDetail(Req)
	if err != nil {
		beego.Error(constant.ErrSystem, err)
		this.ResponseInfo(500, constant.ErrSystem, nil)
	}
	this.ResponseInfo(200, nil, res)
}

// TxGraphData 交易图谱-交易图
func (this *TxController) AddressTxGraph() {
	this.IsPost()
	Req := apiModels.ReqAddressTxGraph{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	res, err := this.txService.TxGraphData(Req)
	if err != nil {
		beego.Error(constant.ErrSystem, err)
		this.ResponseInfo(500, constant.ErrSystem, nil)
	}
	this.ResponseInfo(200, nil, res)
}
