package baseClass

import (
	"github.com/astaxie/beego/orm"
)

type CarrierList struct {
	Id      int
	Carrier string //运营商
	Address string //地区
	DnsIp   string `orm:"unique"` //DNS服务器IP
}

func CarrierGet(carrier, address string) (count int64, carriers []CarrierList, err error) {
	qs := orm.NewOrm().QueryTable("carrier_list")
	if carrier != "" {
		qs = qs.Filter("carrier", carrier)
	}
	if address != "" {
		qs = qs.Filter("address", address)
	}
	if count, err = qs.All(&carriers); err != nil {
		return 0, nil, err
	}
	return count, carriers, nil
}

func (c *CarrierList) CarrierAdd() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		//beego.Error("insert manager error : ", err)
		return err
	}
	return nil
}

func (c *CarrierList) CarrierUpdate() (flag int64, err error) {
	var sel_carrier []CarrierList
	if num, err := orm.NewOrm().
		Raw("SELECT * FROM carrier_list WHERE id != ? and dns_ip = ?", c.Id, c.DnsIp).
		QueryRows(&sel_carrier); err == nil && num > 0 {
		return 1, nil //IP已重复
	} else if err != nil {
		return 0, err
	} else {
		if _, err = orm.NewOrm().Update(c); err != nil {
			return 0, err
		}
		return 2, nil //修改成功
	}
}

func CarrierDelete(id int) error {
	qs := orm.NewOrm().QueryTable("carrier_list").Filter("id", id)
	if _, err := qs.Delete(); err != nil {
		return err
	}
	return nil
}
