package controllers

import (
	"github.com/astaxie/beego"
	"github.com/server/txs-analysis/service"
)

type RouteQuery struct {
	BaseController
	addressService service.AddressService
	txService      service.TxService
	nebulaService  service.NebulaService
}

//nebula:根据指定地址，查询所有入账记录
func (this *RouteQuery) EntryTxsByAddress() {
	address := this.GetString("address")
	if len(address) != 42 {
		beego.Error("input address error. address:", address)
		this.ResponseInfo(500, "input param address error.", nil)
		return
	}
	//result, err := this.nebulaService.GetEntryTxsByAddress(address)
	//if err != nil {
	//	beego.Error("GetEntryTxsByAddress error.", err)
	//	this.ResponseInfo(500, "GetEntryTxsByAddress error.", err)
	//	return
	//}
	this.ResponseInfo(200, nil, nil)
}
