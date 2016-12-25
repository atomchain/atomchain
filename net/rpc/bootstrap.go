package rpc

import (
	"chain/common/config"
	_ "chain/net/rpc/routers"
	"fmt"

	"github.com/astaxie/beego"
)

// Bootstrap start rpc server for RESTful API
func Bootstrap() {
	fmt.Println("RPC Framework Start...")
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.RecoverPanic = true
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.AppConfig.Set("httpport", fmt.Sprint(config.Parameters.HttpInfoPort))
	go beego.Run()
}
