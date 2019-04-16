package baseTable

import (
	"CcsuWebMonitor/controllers"

	"encoding/json"

	"CcsuWebMonitor/models/baseClass"

	"github.com/astaxie/beego"
)

type WebSiteController struct {
	controllers.BaseController
	baseClass.WebSiteList
}

type siteSearch struct {
	WebName string
	Url     string
}

//@router /ccsu/base_table/site_get [get]
func (w *WebSiteController) GetWebSite() {
	if controllers.LoginCheck(&w.BaseController) {
		var search siteSearch
		json.Unmarshal(w.Ctx.Input.RequestBody, &search)
		if count, websites, err := baseClass.WebSiteGet(search.WebName, search.Url); err != nil {
			beego.Error("Site get error : ", err)
			w.Data["json"] = controllers.ResultJson(7001, "获取失败")
			w.ServeJSON()
			return
		} else {
			result := controllers.ResultJson(6000, "")
			result["count"] = count
			result["websites"] = websites
			w.Data["json"] = result
			w.ServeJSON()
			return
		}
	}
}

//@router /ccsu/base_table/site_add [post]
func (w *WebSiteController) AddWebSite() {
	if controllers.LoginCheck(&w.BaseController) {
		json.Unmarshal(w.Ctx.Input.RequestBody, &w.WebSiteList)
		switch w.WebSiteAdd() {
		case 0:
			w.Data["json"] = controllers.ResultJson(7000, "添加成功")
			break
		case 1:
			w.Data["json"] = controllers.ResultJson(7000, "网站名已存在")
			break
		case 2:
			w.Data["json"] = controllers.ResultJson(7000, "该链接已存在")
			break
		case 3:
			w.Data["json"] = controllers.ResultJson(7001, "添加失败")
			break
		}
		w.ServeJSON()
		return
	}
}

//@router /ccsu/base_table/site_update [post]
func (w *WebSiteController) UpDateWebSite() {
	if controllers.LoginCheck(&w.BaseController) { //确认是否登陆
		json.Unmarshal(w.Ctx.Input.RequestBody, &w.WebSiteList)

		//check whether getstring is not null
		//if CheckInputString(w.Url, w.WebName) {
		//update mysql success
		/*if !baseClass.WebSiteUpdate(w.Id, w.Url, w.WebName, w.ResponseSleepTime, w.ResolveSleepTime) {
			w.Data["json"] = controllers.ResultJson(7001, "")
			w.ServeJSON()
			return
		}
		w.Data["json"] = controllers.ResultJson(7000, "")
		w.ServeJSON()
		return
		//}
		/*w.Data["json"] = controllers.ResultJson(4001, "")//the getstring have null
		w.ServeJSON()*/
		//return
	}
}

//@router /ccsu/base_table/site_delete [post]
func (w *WebSiteController) DeleteWebSite() {
	if controllers.LoginCheck(&w.BaseController) {
		var urls []string
		json.Unmarshal(w.Ctx.Input.RequestBody, &urls)
		//软删除网站
		for n := range urls {
			if !baseClass.WebSiteDelete(urls[n]) {
				w.Data["json"] = controllers.ResultJson(7001, "删除"+urls[n]+"时失败,操作终止!")
				w.ServeJSON()
				return
			}
		}
		w.Data["json"] = controllers.ResultJson(7000, "")
		w.ServeJSON()
		return
	}
}
