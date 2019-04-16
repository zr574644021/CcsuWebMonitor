package main

import (
	_ "CcsuWebMonitor/initData"
	_ "CcsuWebMonitor/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
