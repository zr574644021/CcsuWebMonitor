package models

import (
	"CcsuWebMonitor/models/baseClass"
	"CcsuWebMonitor/models/resolveClass"
	"CcsuWebMonitor/models/responseClass"
	"CcsuWebMonitor/models/serverClass"

	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(baseClass.CarrierList), new(baseClass.User), new(baseClass.LoginRecord), new(baseClass.WebSiteList))
	orm.RegisterModel(new(baseClass.DnsMonitor), new(baseClass.HttpMonitor))
	orm.RegisterModel(new(resolveClass.DnsVisitError), new(resolveClass.DnsVisitRecord))
	orm.RegisterModel(new(responseClass.UrlVisitError), new(responseClass.UrlVisitRecord))
	orm.RegisterModel(new(serverClass.Server), new(serverClass.Department), new(serverClass.Manager))
}
