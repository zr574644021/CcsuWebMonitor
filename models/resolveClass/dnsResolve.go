package resolveClass

import (
	"CcsuWebMonitor/models/baseClass"
	"CcsuWebMonitor/util"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type ResolveList struct {
	Id          int
	Status      int                    `orm:"default(0)"`
	MonitorTime int64                  `orm:"default(null)"`
	Url         *baseClass.WebSiteList `orm:"rel(fk)"`
}

type DnsVisitRecord struct {
	Id          int
	ErrorMsg    string //异常信息
	ResolveIp   string //解析IP
	MonitorTime int64  `orm:"null"`       //监测时间
	Status      int    `orm:"default(0)"` //状态信息 0:解析正常 1:解析异常 2:解析超时
	ImgStatus   string
	Url         *baseClass.WebSiteList `orm:"rel(fk)"`
	Carrier     *baseClass.CarrierList `orm:"rel(fk)"`
}

type DnsVisitError struct {
	Id       int
	Time     int64                  //异常时间
	ErrorIp  string                 //异常IP
	ErrorMsg string                 //异常信息
	Url      *baseClass.WebSiteList `orm:"rel(fk)"`
	Carrier  *baseClass.CarrierList `orm:"rel(fk)"`
}

//URL解析
func DnsResolve(url, IP string, res chan string) {
	remsg, n, _ := util.Send(IP+":53", url) //地址解析为IP
	if n < 4 {
		close(res)
		return
	}
	ip := strconv.Itoa(int(remsg[n-4])) + "." + strconv.Itoa(int(remsg[n-3])) + "." +
		strconv.Itoa(int(remsg[n-2])) + "." + strconv.Itoa(int(remsg[n-1]))
	res <- ip
	close(res)
}

//获取所有域名解析异常记录
func DnsErrorGet() (int64, []DnsVisitError, error) {
	var dnsError []DnsVisitError
	o := orm.NewOrm()
	i, err := o.QueryTable("dns_visit_error").OrderBy("-time").RelatedSel().All(&dnsError)
	if err != nil {
		return 0, nil, err
	}
	return i, dnsError, nil
}

//获取最新监测记录
func DnsVistGet() (int64, []DnsVisitRecord, error) {
	var dnsVisit []DnsVisitRecord
	i, err := orm.NewOrm().QueryTable("dns_visit_record").OrderBy("dns_id", "carrier_id").RelatedSel().All(&dnsVisit)
	if err != nil {
		return 0, nil, err
	}
	return i, dnsVisit, nil
}

/*//更新最后次检测时间
func (c *DnsList) UpdateTime() {
	c.LastVisitTime = time.Now().Unix()
	if _, err := orm.NewOrm().Update(c); err != nil {
		beego.Error("save update time error is ", err)
		return
	}
}

//保存异常信息
func (c *DnsVisitError) SaveError() bool {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		beego.Error("save error is ", err)
		return false
	}
	return true
}

//判断是否已完成休眠
func (c *DnsList) FinishSleep() bool {
	now := time.Now()
	lasttime := time.Unix(c.LastVisitTime, 0)
	//将两个时间戳相减，判断是否大于等于休眠时间
	if now.Sub(lasttime).Minutes() >= (c.SleepTime - 1) {
		return true
	}
	return false
}

func (c *DnsVisitRecord) SaveVist() bool {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		beego.Error("save error is ", err)
		return false
	}
	return true
}

func DeleteVist() bool {
	_, dnslist, err := DnsGet()
	if err != nil {
		return false
	}
		orm.NewOrm().QueryTable("dns_visit_record").Filter("dns_id", u.Id).Delete()
	}
	return true
}*/
