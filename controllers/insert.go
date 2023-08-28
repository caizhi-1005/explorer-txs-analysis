package controllers

import (
	"github.com/astaxie/beego"
	"github.com/server/txs-analysis/models/nebulaModels"
	"github.com/server/txs-analysis/service"
)

type RouteInsert struct {
	BaseController
	nebulaService service.NebulaService
}

func (this *RouteInsert) InitNebula() {
	beego.Debug("InitNebula start------------>")
	nebulaDB := nebulaModels.Init()
	//准备tag
	err := nebulaModels.PrepareAddress(nebulaDB)
	if err != nil {
		beego.Error("nebulaModels.PrepareAddress error: ", err)
		this.ResponseInfo(500, "PrepareAddress error.", err)
		return
	}
	//准备edge
	err = nebulaModels.PrepareTxs(nebulaDB)
	if err != nil {
		beego.Error("nebulaModels.PrepareTxs error: ", err)
		this.ResponseInfo(500, "PrepareTxs error.", err)
		return
	}
	nebulaDB.Close()
	beego.Info("nebula prepare success")
	this.ResponseInfo(200, "nebula prepare success succeed!", "ok")
}
