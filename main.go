package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/server/txs-analysis/routers"
)

func main() {
	initLogger()
	beego.Info(beego.BConfig.AppName, "started!")
	beego.Run()
}

func initLogger() (err error) {

	config := make(map[string]interface{})
	config["filename"] = beego.AppConfig.String("logPath")

	// map 转 json
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger failed, marshal err:", err)
		return
	}
	// log配置
	beego.SetLogger(logs.AdapterFile, string(configStr))
	// log打印文件名和行数
	beego.SetLogFuncCall(true)
	return
}
