package baseClass

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type WebSiteList struct {
	Id      int
	Url     string `orm:"unique"` //网站链接
	WebName string `orm:"unique"` //网站名称
	Ip      string //网站IP
	Status  int    `orm:"default(0)"` //软删除 0:未删除 1:已删除
	/*HttpStatus        int    `orm:"default(0)"` //(响应)状态码 0:不检测  1:检测
	DnsStatus         int    `orm:"default(0)"` //(解析)状态码
	ResponseSleepTime int64  `orm:"null"`       //响应监测频率
	ResolveSleepTime  int64  `orm:"null"`       //解析频率
	ResponseTime      int64  `orm:"null"`       //响应时间
	ResolveTime       int64  `orm:"null"`       //解析时间*/
}

//存储需响应网站
type HttpMonitor struct {
	Id                int
	ResponseSleepTime int
	ResponseTime      int
	Qsrq              int
	WebSite           *WebSiteList `orm:"rel(fk)"`
}

//存储需解析网站
type DnsMonitor struct {
	Id               int
	ResolveSleepTime int
	ResolveTime      int
	Qsrq             int
	WebSite          *WebSiteList `orm:"rel(fk)"`
}

//获取所有需要解析的域名记录
func WebSiteGet(name, url string) (count int64, webSites []WebSiteList, err error) {
	var WebSiteList []WebSiteList
	qs := orm.NewOrm().QueryTable("web_site_list").Filter("status", 0) //筛选删除标志为0的记录
	if name != "" {
		qs = qs.Filter("web_name", name)
	}
	if url != "" {
		qs = qs.Filter("url", url)
	}
	if count, err = qs.All(&webSites); err != nil {
		return 0, nil, err
	}
	return count, WebSiteList, nil
}

func (w *WebSiteList) WebSiteAdd() int {
	o := orm.NewOrm().QueryTable("web_site_list").Filter("status", 0)
	qs := orm.NewOrm()
	var webSite WebSiteList
	if o.Filter("web_name", w.WebName).One(&webSite) == nil {
		return 1 //网站名已存在
	}
	if o.Filter("url", w.Url).One(&webSite) == nil {
		return 2 //链接已存在
	}
	if o.Filter("url", w.Url).Filter("web_name", w.WebName).Filter("status", 1).One(&webSite) == nil {
		webSite.Status = 0
		if _, err := orm.NewOrm().Update(&webSite); err != nil {
			beego.Error("web_url_add update url_list_status error is ", err)
			return 3 //添加失败
		}
		return 0 //添加成功
	}
	if _, err := qs.Insert(w); err != nil {
		beego.Error("web url add error is ", err)
		return 3
	}
	return 0
}

// 更新网站
func (w *WebSiteList) WebSiteUpdate(id int, url, webName string, responseSleepTime, resolveSleepTime int64) bool {
	var webSite WebSiteList
	o := orm.NewOrm()
	if err := o.QueryTable("web_site_list").
		Filter("id", id).RelatedSel().One(&webSite); err != nil {
		beego.Error("web_url_update get url_list error is ", err)
		return false
	} else {
		if url != "" {
			webSite.Url = url
		}
		if webName != "" {
			webSite.WebName = webName
		}
		if _, err := o.Update(&webSite); err != nil {
			beego.Error("web_url_update update url_list error is ", err)
			return false
		}
	}
	return true
}

func WebSiteDelete(url string) bool {
	var webSite WebSiteList
	o := orm.NewOrm()
	if err := o.QueryTable("web_site_list").
		Filter("url", url).Filter("status", 0).One(&webSite); err != nil {
		beego.Error("web_url_delete get url_list error is ", err)
		return false
	} else {
		webSite.Status = 1
		// 软删除，Delete属性为删除的时候为1
		if _, err := o.Update(&webSite); err != nil {
			beego.Error("web_url_delete update url_list error is ", err)
			return false
		}
	}
	return true
}
