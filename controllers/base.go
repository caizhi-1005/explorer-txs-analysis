package controllers

import (
	"github.com/astaxie/beego"
	"github.com/server/txs-analysis/constant"
	"github.com/server/txs-analysis/utils"
)

type BaseController struct {
	beego.Controller
}

// 响应
func (d *BaseController) ResponseInfo(code int, errMsg interface{}, result interface{}) {
	switch code {
	case 200:
		d.Data["json"] = map[string]interface{}{"code": 200, "err_msg": errMsg, "data": result}
	default:
		d.Data["json"] = map[string]interface{}{"code": 500, "err_msg": errMsg, "data": result}
	}
	d.ServeJSON()
}

// 校验POST请求
func (this *BaseController) IsPost() bool {
	if this.Ctx.Request.Method != "POST" {
		beego.Error("Request method is not post.")
		this.ResponseInfo(500, constant.ErrRequestMode, nil)
	}
	return true
}

// 分页
func (this *BaseController) Pagination(ReqPage, ReqLength string) (page, length int64) {
	defaultPage := beego.AppConfig.DefaultInt64("pagination::page", 1)
	defaultLength := beego.AppConfig.DefaultInt64("pagination::length", 10)
	page = utils.StringInt64(ReqPage)
	length = utils.StringInt64(ReqLength)

	if page <= 0 {
		page = defaultPage
	}
	if length <= 0 {
		length = defaultLength
	}
	return page, length
}
