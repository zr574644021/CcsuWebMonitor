package responseClass

import (
	"CcsuWebMonitor/models/baseClass"

	"github.com/astaxie/beego/orm"
)

type UrlVisitRecord struct {
	Id           int
	ResponseTime float64                //响应时间
	Time         int64                  //监测时间
	Url          *baseClass.WebSiteList `orm:"rel(fk)"`
}

type UrlVisitError struct {
	Id          int
	Time        int64                  //监测时间
	ErrorStatus string                 //状态码
	Url         *baseClass.WebSiteList `orm:"rel(fk)"`
}

type StatusTime struct {
	Status string
	Time   float64
}

/*// 获取所有需要监测的记录
func UrlGet() (int64, []baseClass.UrlList, error) {
	var urllists []baseClass.UrlList
	i, err := orm.NewOrm().QueryTable("url_list").
		Filter("status", 0).All(&urllists)
	//筛选只有0的，1为软删除
	if err != nil {
		return 0, nil, err
	}
	return i, urllists, nil
}*/

// 获取最近的错误记录,can select limit max number
func UrlErrorGet() (int64, *[]UrlVisitError, error) {
	var urlerror []UrlVisitError
	o := orm.NewOrm()
	i, err := o.QueryTable("url_visit_error").OrderBy("-time").RelatedSel().All(&urlerror)
	if err != nil {
		return 0, nil, err
	}

	return i, &urlerror, nil
}

func UrlRecordGet(url string, begin, end int64) *[]UrlVisitRecord {
	var urlrecord []UrlVisitRecord
	var urllist baseClass.WebSiteList
	o := orm.NewOrm()
	err1 := o.QueryTable("url_list").
		Filter("status", 0).
		Filter("url", url).
		One(&urllist)
	if err1 != nil {
		return nil
	}
	_, err2 := o.QueryTable("url_visit_record").
		Filter("url_id", &urllist).
		Filter("time__gte", begin).
		Filter("time__lte", end).
		All(&urlrecord)
	if err2 != nil {
		return nil
	}
	return &urlrecord
}

/*// 判断是否已完成休眠
func (c *UrlList) FinishSleep() bool {
	now := time.Now()
	// 将int64的数据转化为时间戳
	lasttime := time.Unix(c.LastVisitTime, 0)
	// 将两个时间戳相减，判断是否大于等于休眠时间
	if now.Sub(lasttime).Minutes() >= (c.SleepTime - 1) {
		return true
	}
	return false
}

// 将监测的记录存入url_list表中(好像没有使用)
func (c *UrlVisitRecord) VisitSave() bool {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		beego.Error("save visit record is ", err)
		return false
	}
	return true
}

func (c *baseClass.User) TimeUpdate() {
	c.LastVisitTime = time.Now().Unix()
	if _, err := orm.NewOrm().Update(c); err != nil {
		beego.Error("save update time error is ", err)
		return
	}
}

// 保存超时监测记录
func (c *UrlVisitError) VisitSave() {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		beego.Error("save visit record is ", err)
		return
	}
	return
}

// 通过URL获取对应网站的请求的休眠时间(-1为查询错误)
func SleepTimeGet(url string) (float64, error) {
	var urllist baseClass.UrlList
	err := orm.NewOrm().QueryTable("url_list").
		Filter("status", 0).Filter("url", url).
		One(&urllist, "sleep_time")
	if err != nil {
		return -1, err
	}
	return urllist.SleepTime, nil
}*/

func IntervalOneRecordGet(url string, begin, end int64) *UrlVisitRecord {
	var urllistrecord UrlVisitRecord
	var urllist baseClass.WebSiteList
	o := orm.NewOrm()
	err1 := o.QueryTable("url_list").
		Filter("url", url).One(&urllist)
	if err1 != nil {
		return nil
	}
	err2 := o.QueryTable("url_visit_record").
		Filter("url_id", &urllist).Filter("time__gte", begin).Filter("time__lte", end).One(&urllistrecord, "response_time", "time")
	if err2 != nil {
		return nil
	}
	return &urllistrecord
}

func UrlErrorRecordQuery(url string, begin, end int64) (int64, *[]UrlVisitError) {
	var urllisterror []UrlVisitError
	var urllist baseClass.WebSiteList
	o := orm.NewOrm()
	err1 := o.QueryTable("url_list").
		Filter("url", url).One(&urllist)
	if err1 != nil {
		return 0, nil
	}
	num, err2 := o.QueryTable("url_visit_record").
		Filter("url_id", &urllist).
		Filter("time__gte", begin).
		Filter("time__lte", end).
		All(&urllisterror)
	if err2 != nil {
		return 0, nil
	}
	return num, &urllisterror
}

// delete url record
func UrlErrorDelete(url string) (int64, error) {
	var urllist baseClass.WebSiteList
	var num int64
	err := orm.NewOrm().
		QueryTable("url_list").
		Filter("url", url).
		One(&urllist)
	if err != nil {
		return 0, err
	}
	var urlrecord UrlVisitRecord
	urlrecord.Url = &urllist
	num, err = orm.NewOrm().
		QueryTable("url_visit_error").
		Filter("url_id", &urllist).Delete()
	if err != nil {
		return 0, err
	}
	return num, nil
}
