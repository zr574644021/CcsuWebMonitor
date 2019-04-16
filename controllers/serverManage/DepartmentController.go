package serverManage

import (
	"CcsuWebMonitor/controllers"

	"encoding/json"

	"CcsuWebMonitor/models/common"
	"CcsuWebMonitor/models/serverClass"

	"github.com/astaxie/beego"
)

type DmController struct {
	controllers.BaseController
	department serverClass.Department
}

type departmentSet struct {
	Id      int
	Name    string
	NewName string
}

//@router /ccsu/server_manage/department_get [*]
func (s *DmController) GetDepartment() {
	if controllers.LoginCheck(&s.BaseController) {
		var page common.Paging
		json.Unmarshal(s.Ctx.Input.RequestBody, &page)
		if count, departments, boolean := serverClass.DepGet(page.Page, page.PageSize); boolean == false {
			result := controllers.ResultJson(7001, "获取失败")
			s.Data["json"] = result
		} else {
			result := controllers.ResultJson(6000, "")
			result["count"] = count
			result["departments"] = departments
			s.Data["json"] = result

		}
		s.ServeJSON()
	}
}

//@router /ccsu/server_manage/department_add [*]
func (s *DmController) AddDepartment() {
	if controllers.LoginCheck(&s.BaseController) {
		json.Unmarshal(s.Ctx.Input.RequestBody, &s.department)
		flag := s.department.DepAdd()
		switch flag {
		case 0:
			s.Data["json"] = controllers.ResultJson(7000, "添加成功")
			break
		case 1:
			s.Data["json"] = controllers.ResultJson(7001, "添加异常")
			break
		case 2:
			s.Data["json"] = controllers.ResultJson(7001, "部门"+s.department.Name+"已存在")
			break
		}
		s.ServeJSON()
	}
}

//@router ccsu/server_manage/department_delete [*]
func (s *DmController) RemoveDepartment() {
	if controllers.LoginCheck(&s.BaseController) {
		json.Unmarshal(s.Ctx.Input.RequestBody, &s.department)
		if count, _ := serverClass.ManagerCheck(s.department.Id, ""); count > 0 {
			s.Data["json"] = controllers.ResultJson(7001, "请先确认该部门下无所属管理员")
		} else {
			if boolean := s.department.DepDelete(); boolean == false {
				s.Data["json"] = controllers.ResultJson(7001, "删除失败")
			} else {
				s.Data["json"] = controllers.ResultJson(7000, "删除成功")
			}
		}
		s.ServeJSON()
	}
}

//@router /ccsu/server_manage/department_update [*]
func (s *DmController) SetDepartment() {
	if controllers.LoginCheck(&s.BaseController) {
		var department serverClass.Department
		var departmentUp departmentSet
		json.Unmarshal(s.Ctx.Input.RequestBody, &departmentUp)
		// use the name to set data
		department.Name = departmentUp.NewName
		newName := departmentUp.NewName
		if count, _ := department.DepQuery(department.Name); count > 0 {
			s.Data["json"] = controllers.ResultJson(7001, "该部门已存在")
		} else {
			department.Id = departmentUp.Id
			department.Name = departmentUp.Name
			if err := department.DepUpdate(newName); err != nil {
				beego.Error("Update department error : ", err)
				s.Data["json"] = controllers.ResultJson(7001, "更新失败")
			} else {
				s.Data["json"] = controllers.ResultJson(7000, "更新成功")
			}
		}
		s.ServeJSON()
	}
}
