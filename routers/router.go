package routers

import (
	"CcsuWebMonitor/controllers"

	"CcsuWebMonitor/controllers/serverManage"

	"CcsuWebMonitor/controllers/webMonitor"

	"CcsuWebMonitor/controllers/baseTable"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.Include(&controllers.MainController{})
	beego.Include(&controllers.BaseController{})
	beego.Include(&serverManage.ServerController{})
	beego.Include(&serverManage.DmController{})
	beego.Include(&serverManage.ManagerController{})
	beego.Include(&webMonitor.HttpController{})
	beego.Include(&baseTable.CarrierController{})

	//解决跨域问题
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
}
