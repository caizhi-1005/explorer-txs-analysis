package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/server/txs-analysis/constant"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/service"
)

type NFTController struct {
	BaseController
	nftService service.NftService
	nebulaService  service.NebulaService
}

// NFTList NFT溯源-全部NFT列表(下拉列表)
func (this *NFTController) NFTList() {
	this.IsPost()
	result, err := this.nftService.NFTList()
	if err != nil {
		beego.Error("NFT NFTList error.", err)
		this.ResponseInfo(500, "NFT NFTList error.", err)
		return
	}
	this.ResponseInfo(200, nil, result)
}

// TokenIdList NFT溯源-Token ID下拉列表
func (this *NFTController) TokenIdList() {
	this.IsPost()
	Req := apiModels.ReqTokenIDList{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}
	result, err := this.nftService.TokenIDList(Req)
	if err != nil && err != orm.ErrNoRows {
		beego.Error("NFT TokenIdList error.", err)
		this.ResponseInfo(500, "NFT TokenIdList error.", err)
		return
	}
	this.ResponseInfo(200, nil, result)
}

// NFTAddressDetail NFT溯源-地址详情
func (this *NFTController) NFTAddressDetail() {
	this.IsPost()
	//Req := apiModels.ReqNFTAddressDetail{}
	//if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
	//	beego.Error(constant.ErrParam, err)
	//	this.ResponseInfo(500, constant.ErrParam, nil)
	//}
	//
	////Res := apiModels.RespNFTAddressDetail{}
	//result, err := this.nftService.NFTAddressDetail(Req)
	//if err != nil && err != orm.ErrNoRows{
	//	beego.Error("NFT NFTAddressDetail error.", err)
	//	this.ResponseInfo(500, "NFT NFTAddressDetail error.", err)
	//	return
	//}
	//this.ResponseInfo(200, nil, result)
}

// NFTTransferDetailByAddress NFT溯源-地址详情-流转详情
func (this *NFTController) NFTTransferDetailByAddress() {
	this.IsPost()
	Req := apiModels.ReqNFTTransferDetailsByAddress{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}

	page, pageSize := this.Pagination(Req.Start, Req.Length)
	Req.Page = page
	Req.PageSize = pageSize

	result, err := this.nftService.NFTTransferDetailByAddress(Req)
	if err != nil && err != orm.ErrNoRows{
		beego.Error("NFT NFTTransferDetailByAddress error.", err)
		this.ResponseInfo(500, "NFT NFTTransferDetailByAddress error.", err)
		return
	}
	this.ResponseInfo(200, nil, result)
}

// NFTTxDetail NFT溯源-交易详情
func (this *NFTController) NFTTxDetail() {
	this.IsPost()
	Req := apiModels.ReqNFTTxDetail{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}

	result, err := this.nftService.NFTTxDetail(Req)
	if err != nil {
		beego.Error("NFT NFTTxDetail error.", err)
		this.ResponseInfo(500, "NFT NFTTxDetail error.", err)
		return
	}
	this.ResponseInfo(200, nil, result)
}

// NFTDetail NFT溯源-NFT详情
func (this *NFTController) NFTDetail() {
	this.IsPost()
	Req := apiModels.ReqNFTDetail{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}

	result, err := this.nftService.NFTDetail(Req)
	if err != nil {
		beego.Error("NFT NFTDetail error.", err)
		this.ResponseInfo(500, "NFT NFTDetail error.", err)
		return
	}
	this.ResponseInfo(200, nil, result)
}

// NFTTransferDetailByTokenId NFT溯源-NFT详情-NFT流转详情列表
func (this *NFTController) NFTTransferDetailByTokenId() {
	this.IsPost()
	Req := apiModels.ReqNFTTransferDetailsByTokenId{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}

	result, err := this.nftService.NFTTransferDetailByTokenId(Req)
	if err != nil {
		beego.Error("NFT NFTTransferDetailByTokenId error.", err)
		this.ResponseInfo(500, "NFT NFTTransferDetailByTokenId error.", err)
		return
	}
	this.ResponseInfo(200, nil, result)
}


// NFTStartAnalysis NFT溯源-NFT开始分析-交易图
//todo 按照前端需求，返回指定格式
func (this *NFTController) NFTStartAnalysis() {
	this.IsPost()
	Req := apiModels.ReqNFTStartAnalysis{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}

	result, err := this.nebulaService.NFTStartAnalysis(Req.ContractAddress, Req.Input)
	if err != nil {
		beego.Error("NFT NFTStartAnalysis error.", err)
		this.ResponseInfo(500, "NFT NFTStartAnalysis error.", err)
		return
	}
	this.ResponseInfo(200, nil, result)
}


// NFTTrace NFT溯源-NFT追溯-交易图
//todo 按照前端需求，返回指定格式
func (this *NFTController) NFTTrace() {
	this.IsPost()
	Req := apiModels.ReqNFTTrace{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &Req); nil != err {
		beego.Error(constant.ErrParam, err)
		this.ResponseInfo(500, constant.ErrParam, nil)
	}

	result, err := this.nebulaService.GetNFTTxsPath(Req)
	if err != nil {
		beego.Error("NFT NFTTransferDetailByTokenId error.", err)
		this.ResponseInfo(500, "NFT NFTTransferDetailByTokenId error.", err)
		return
	}
	this.ResponseInfo(200, nil, result)
}
