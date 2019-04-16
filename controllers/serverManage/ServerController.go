package serverManage

import (
	"CcsuWebMonitor/controllers"
	"encoding/json"

	"CcsuWebMonitor/models/serverClass"
)

type ServerController struct {
	controllers.BaseController
	server serverClass.Server
}

//@router /ccsu/server_manage/server_get [post]
func (s *ServerController) GetServer() {
	if controllers.LoginCheck(&s.BaseController) {
		var search serverClass.ServerSearch
		json.Unmarshal(s.Ctx.Input.RequestBody, &search)
		count, servers, boolean := serverClass.ServerGet(search)
		if boolean == false {
			result := controllers.ResultJson(7001, "获取失败")
			s.Data["json"] = result
		} else {
			result := controllers.ResultJson(6000, "")
			result["servers"] = servers
			result["count"] = count
			s.Data["json"] = result
		}
		s.ServeJSON()
	}
}

//@router /ccsu/server_manage/server_add [post]
func (s *ServerController) AddServer() {
	if controllers.LoginCheck(&s.BaseController) {
		json.Unmarshal(s.Ctx.Input.RequestBody, &s.server)
		flag := s.server.ServerAdd()
		switch flag {
		case 0:
			s.Data["json"] = controllers.ResultJson(7000, "添加成功")
			break
		case 1:
			s.Data["json"] = controllers.ResultJson(7001, "位置"+s.server.Area+"已占用")
			break
		case 2:
			s.Data["json"] = controllers.ResultJson(7001, "添加异常")
			break
		}
		s.ServeJSON()
	}
}

//@router /ccsu/server_manage/server_delete [post]
func (s *ServerController) RemoveServer() {
	if controllers.LoginCheck(&s.BaseController) {
		var server []serverClass.Server
		json.Unmarshal(s.Ctx.Input.RequestBody, &server)
		var i int
		for i = 0; i < len(server); i++ {
			if boolean := serverClass.ServerDelete(server[i].Id); boolean == false {
				s.Data["json"] = controllers.ResultJson(7001, "删除服务器"+server[i].Name+"失败,操作终止")
				break
			}
		}
		if i == len(server) {
			s.Data["json"] = controllers.ResultJson(7000, "删除成功")
		}
		s.ServeJSON()
	}
}

//@router /ccsu/server_manage/server_update [post]
func (s *ServerController) SetServer() {
	if controllers.LoginCheck(&s.BaseController) {
		json.Unmarshal(s.Ctx.Input.RequestBody, &s.server)
		flag := s.server.ServerUpdate()
		switch flag {
		case 0:
			s.Data["json"] = controllers.ResultJson(7000, "更新成功")
			break
		case 1:
			s.Data["json"] = controllers.ResultJson(7001, "位置"+s.server.Area+"已占用")
			break
		case 2:
			s.Data["json"] = controllers.ResultJson(7001, "更新异常")
			break
		}
		s.ServeJSON()
	}
}
