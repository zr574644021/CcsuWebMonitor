package serverManage

import (
	"CcsuWebMonitor/controllers"
	"CcsuWebMonitor/models/serverClass"

	"encoding/json"

	"CcsuWebMonitor/models/common"
)

type ManagerController struct {
	controllers.BaseController
	manager serverClass.Manager
}

type managerSearch struct {
	Name  string
	DepId int
	Page  common.Paging
}

//@router /ccsu/server_manage/manager_get [post]
func (s *ManagerController) GetManager() {
	if controllers.LoginCheck(&s.BaseController) {
		var search managerSearch
		json.Unmarshal(s.Ctx.Input.RequestBody, &search)
		if count, managers, boolean := serverClass.ManagerGet(search.Name, search.DepId, search.Page.Page, search.Page.PageSize); boolean == false {
			s.Data["json"] = controllers.ResultJson(7001, "获取失败")
		} else {
			result := controllers.ResultJson(6000, "")
			result["count"] = count
			result["managers"] = managers
			s.Data["json"] = result
		}
		s.ServeJSON()
	}
}

//@router /ccsu/server_manage/manager_add [post]
func (s *ManagerController) AddManager() {
	if controllers.LoginCheck(&s.BaseController) {
		json.Unmarshal(s.Ctx.Input.RequestBody, &s.manager)
		flag := s.manager.ManagerAdd()
		switch flag {
		case 0:
			s.Data["json"] = controllers.ResultJson(7000, "添加成功")
			break
		case 1:
			s.Data["json"] = controllers.ResultJson(7001, "添加失败")
			break
		case 2:
			s.Data["json"] = controllers.ResultJson(7001, "添加失败")
			break
		}
		s.ServeJSON()
	}
}

//@router /ccsu/server_manage/manager_delete  [post]
func (s *ManagerController) RemoveManager() {
	if controllers.LoginCheck(&s.BaseController) {
		var manager []serverClass.Manager
		json.Unmarshal(s.Ctx.Input.RequestBody, &manager)
		flag := true
		for i := 0; i < len(manager); i++ {
			if count, boolean := serverClass.ServerCheck(manager[0].Id, ""); count > 0 && boolean == true {
				s.Data["json"] = controllers.ResultJson(7001, "请确认管理员"+manager[i].Name+"下无负责的服务器")
				flag = false
				break
			}
			if flag := serverClass.ManagerDelete(manager[i].Id); !flag {
				s.Data["json"] = controllers.ResultJson(7001, "删除"+manager[i].Name+"时失败,操作终止")
				flag = false
				break
			}
		}
		if flag {
			s.Data["json"] = controllers.ResultJson(7000, "删除成功")
		}
		s.ServeJSON()
	}
}

//@router /ccsu/server_manage/manager_update [post]
func (s *ManagerController) SetManager() {
	if controllers.LoginCheck(&s.BaseController) {
		//var need_set serverClass.Manager
		json.Unmarshal(s.Ctx.Input.RequestBody, &s.manager)
		/*need_set.Name = s.manager.Name
		need_set.PhoneNumber = s.manager.PhoneNumber
		need_set.Dep = s.manager.Dep*/
		flag := s.manager.ManagerUpdate()
		switch flag {
		case 0:
			s.Data["json"] = controllers.ResultJson(7000, "更新成功")
			break
		case 1:
			s.Data["json"] = controllers.ResultJson(7001, s.manager.Name+"已存在")
			break
		case 2:
			s.Data["json"] = controllers.ResultJson(7001, "更新失败")
			break
		}
		s.ServeJSON()
	}
}
