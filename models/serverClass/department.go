package serverClass

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Department struct {
	Id   int
	Name string `orm:"unique"` //部门名
}

//添加部门
func (d *Department) DepAdd() int {
	o := orm.NewOrm()
	count, err := o.QueryTable("department").Filter("name", d.Name).Count()
	if err != nil {
		beego.Error("department insert with query error : ", err)
		return 1
	}
	if count > 0 {
		return 2 //部门已存在
	}
	if _, err := o.Insert(d); err != nil {
		beego.Error("department insert error : ", err)
		return 1
	}
	return 0
}

//删除部门
func (d *Department) DepDelete() bool {
	qs := orm.NewOrm().QueryTable("department")
	if d.Id != 0 {
		qs = qs.Filter("id", d.Id)
	}
	if d.Name != "" {
		qs = qs.Filter("name", d.Name)
	}
	if _, err := qs.Delete(); err != nil {
		beego.Error("department delete error : ", err)
		return false
	}
	return true
}

//查询部门
func (d *Department) DepQuery(name string) (int, bool) {
	//var department Department
	if err := orm.NewOrm().QueryTable("department").Filter("name", name).One(d); err != nil && err != orm.ErrNoRows {
		beego.Error("department query error : ", err)
		return 0, false
	}
	if d.Id != 0 {
		return 1, true
	}
	return 0, true
}

//修改部门
func (d *Department) DepUpdate(newName string) error {
	o := orm.NewOrm()
	if err := o.QueryTable("department").Filter("id", d.Id).One(d); err != nil {
		return err
	} else {
		d.Name = newName
		if _, err := o.Update(d); err != nil {
			return err
		}
	}
	return nil
}

func DepGet(page, pageSize int) (count int64, dps []Department, bool bool) {
	qs := orm.NewOrm().QueryTable("department")
	_, err := qs.Limit(pageSize, (page-1)*pageSize).All(&dps)
	if err != nil {
		beego.Error("department get error : ", err)
		return 0, nil, false
	}
	count, _ = qs.Count()
	return count, dps, true
}
