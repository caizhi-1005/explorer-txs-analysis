package routers

import (
	"github.com/astaxie/beego"
	"github.com/server/txs-analysis/controllers"
)

func init() {
	//初始化nebula (tag:address, edge:transaction)
	beego.Router("/api/insert/init", &controllers.RouteInsert{}, "get:InitNebula")

	//api接口
	ns := beego.NewNamespace("/api",
		//地址分析
		beego.NSNamespace("address",
			//可查的token下拉列表
			beego.NSRouter("/tokenList", &controllers.AddressController{}, "post:AddressTokenList"),
			//地址详情
			beego.NSRouter("/detail", &controllers.AddressController{}, "post:AddressDetail"),
			//地址分析图
			beego.NSRouter("/txAnalysis", &controllers.AddressController{}, "post:AddressTxAnalysis"),
			//todo 前端可直接通过交易图谱获取数据
			//交易详情
			beego.NSRouter("/txDetail", &controllers.AddressController{}, "post:AddressTxDetail"),
			//交易详情-交易列表（分页）
			beego.NSRouter("/txList", &controllers.AddressController{}, "post:AddressTxList"),
		),
		//交易图谱
		beego.NSNamespace("txs",
			//交易详情
			beego.NSRouter("/detail", &controllers.TxController{}, "post:TxDetail"),
			//地址详情
			beego.NSRouter("/addressDetail", &controllers.TxController{}, "post:TxAddressDetail"),
			//交易图
			beego.NSRouter("/txGraph", &controllers.TxController{}, "post:AddressTxGraph"),
		),
		//NFT溯源
		beego.NSNamespace("nft",
			//全部NFT列表(下拉列表)
			beego.NSRouter("/list", &controllers.NFTController{}, "post:NFTList"),
			//Token ID下拉列表
			beego.NSRouter("/tokenIdList", &controllers.NFTController{}, "post:TokenIdList"),
			//地址详情
			beego.NSRouter("/addressDetail", &controllers.NFTController{}, "post:NFTAddressDetail"),
			//地址详情-流转详情列表（分页）
			beego.NSRouter("/transferDetailByAddress", &controllers.NFTController{}, "post:NFTTransferDetailByAddress"),
			//交易详情
			beego.NSRouter("/txDetail", &controllers.NFTController{}, "post:NFTTxDetail"),
			//NFT详情
			beego.NSRouter("/detail", &controllers.NFTController{}, "post:NFTDetail"),
			//NFT详情-NFT流转详情列表（分页）
			beego.NSRouter("/transferDetailByTokenId", &controllers.NFTController{}, "post:NFTTransferDetailByTokenId"),
			//图-开始分析
			beego.NSRouter("/startAnalysis", &controllers.NFTController{}, "post:NFTStartAnalysis"),
			//图-NFT追溯-交易图
			beego.NSRouter("/trace", &controllers.NFTController{}, "post:NFTTrace"),
		),
	)
	beego.AddNamespace(ns)

}
