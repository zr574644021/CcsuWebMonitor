package serverClass

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Manager struct {
	Id          int
	Name        string `orm:"unique"` // 管理员姓名
	PhoneNumber string
	Dep         *Department `orm:"rel(fk)" ` //部门外键
	Delete      bool        `orm:"default(false)"`
}

//新增管理员
func (m *Manager) ManagerAdd() int {
	if count, _ := ManagerCheck(0, m.Name); count > 0 {
		return 1
	}
	if _, err := orm.NewOrm().Insert(m); err != nil {
		beego.Error("insert manager error : ", err)
		return 2
	}
	return 0
}

//删除管理员
func ManagerDelete(id int) bool {
	//var num int64
	qs := orm.NewOrm().QueryTable("manager").Filter("id", id)
	if _, err := qs.Count(); err != nil {
		beego.Error("count manager : ", err)
		return false
	}
	if _, err := qs.Delete(); err != nil {
		beego.Error("Delete manager : ", err)
		return false
	}
	return true
}

//管理员检查
func ManagerCheck(depId int, name string) (int64, bool) {
	o := orm.NewOrm().QueryTable("manager")
	if depId != 0 {
		o = o.Filter("dep_id", depId)
	}
	if name != "" {
		o = o.Filter("name", name)
	}
	count, err := o.Count()
	if err != nil {
		beego.Error("manager query error : ", err)
		return 0, false
	}
	return count, true
}

//修改管理员
func (m *Manager) ManagerUpdate() int {
	/*qs := orm.NewOrm().QueryTable("manager").Filter("id", m.Id)
	if err := qs.One(m); err != nil {
		beego.Error("manager update with query error : ", err)
		return false
	}
	if need_update_m.Dep != nil && need_update_m.Dep.Id != m.Dep.Id {
		m.Dep = need_update_m.Dep
	}
	if need_update_m.Name != "" && need_update_m.Name != m.Name {
		m.Name = need_update_m.Name
	}
	if need_update_m.PhoneNumber != "" && need_update_m.PhoneNumber != m.PhoneNumber {
		m.PhoneNumber = need_update_m.PhoneNumber
	}*/
	if count, _ := ManagerCheck(0, m.Name); count > 0 {
		return 1
	}
	if _, err := orm.NewOrm().Update(m); err != nil {
		beego.Error("manager update error : ", err)
		return 2
	}
	return 0
}

func ManagerGet(name string, depId, page, pageSize int) (count int64, managers []Manager, bool bool) {
	qs := orm.NewOrm().QueryTable("manager").Filter("delete", 0).RelatedSel().OrderBy("dep_id", "id")

	if name != "" {
		qs = qs.Filter("name", name)
	}
	if depId != 0 {
		qs = qs.Filter("dep_id", depId)
	}
	_, err := qs.Limit(pageSize, (page-1)*pageSize).All(&managers)
	if err != nil {
		beego.Error("manager query error : ", err)
		bool = false
		count = 0
	}
	bool = true
	count, _ = qs.Count()
	return count, managers, bool
}
