package serverClass

import (
	"CcsuWebMonitor/models/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Server struct {
	Id           int
	Area         string   // 区域
	HardPosition string   // 物理位置
	Name         string   // 服务器名
	Message      string   // 服务器介绍
	HardWare     string   // 硬件环境
	Manager      *Manager `orm:"rel(fk)"` // 管理员
}

type lsServer struct {
	Id           int
	Area         string
	HardPosition string
	Name         string
	Message      string
	HardWare     string
	MaId         int
	MaName       string
	MaPhone      string
	DepId        int
	DepName      string
}

type ServerSearch struct {
	Area    string
	MaName  string
	DepName string
	Page    common.Paging
}

//新增服务器
func (s *Server) ServerAdd() int {
	o := orm.NewOrm()
	if count, _ := ServerCheck(0, s.Area); count > 0 {
		return 2
	}
	if _, err := o.Insert(s); err != nil {
		beego.Error("server insert error : ", err)
		return 1 //
	}
	return 0
}

//删除服务器
func ServerDelete(id int) bool {
	qs := orm.NewOrm().QueryTable("server").Filter("id", id)
	if _, err := qs.Delete(); err != nil {
		beego.Error("server delete error : ", err)
		return false
	}
	return true
}

//修改服务器信息
func (s *Server) ServerUpdate() int {
	o := orm.NewOrm()
	var err error
	var num int
	if err = o.Raw("select count(id) as num from server where id != ? and hard_position = '"+s.HardPosition+"'", s.Id).QueryRow(&num); err != nil && err != orm.ErrNoRows {
		beego.Error("server count select by area error : ", err)
		return 2
	}
	if num > 0 {
		return 1 //位置被占用
	}
	var server Server
	if err = o.QueryTable("server").
		Filter("id", s.Id).RelatedSel().One(&server); err != nil {
		beego.Error("server update with query error : ", err)
		return 2
	}
	if s.Name != "" {
		server.Name = s.Name
	}
	if s.HardPosition != "" {
		server.HardPosition = s.HardPosition
	}
	if s.Area != "" {
		server.Area = s.Area
	}
	if s.HardWare != "" {
		server.HardWare = s.HardWare
	}
	if s.Manager != nil {
		server.Manager.Id = s.Manager.Id
	}
	if _, err = o.Update(&server); err != nil {
		beego.Error("server update error : ", err)
		return 2
	}
	return 0
}

func ServerGet(search ServerSearch) (int64, []Server, bool) {
	qs := orm.NewOrm().Raw("SELECT a.id,a.area,a.hard_position,a.name,a.message,a.hard_ware,"+
		"b.id as ma_id,b.name as ma_name,b.phone_number as ma_phone,c.id as dep_id,c.name as dep_name FROM server a "+
		"INNER JOIN manager b ON a.manager_id=b.id "+
		"INNER JOIN department c ON c.id=b.dep_id "+
		"WHERE a.area LIKE '%"+search.Area+"%' AND b.name LIKE '%"+search.MaName+"%' AND c.name LIKE '%"+search.DepName+"%' "+
		"LIMIT ? OFFSET ?",
		search.Page.PageSize, (search.Page.Page-1)*search.Page.PageSize)
	var lsServers []lsServer
	if _, err := qs.QueryRows(&lsServers); err != nil {
		beego.Error("server get error : ", err)
		return 0, nil, false
	}
	servers := make([]Server, len(lsServers))
	for i, a := range lsServers {
		servers[i].Manager = new(Manager)
		servers[i].Manager.Dep = new(Department)
		servers[i].Id = a.Id
		servers[i].Area = a.Area
		servers[i].Name = a.Name
		servers[i].Message = a.Message
		servers[i].HardPosition = a.HardPosition
		servers[i].HardWare = a.HardWare
		servers[i].Manager.Id = a.MaId
		servers[i].Manager.Name = a.MaName
		servers[i].Manager.PhoneNumber = a.MaPhone
		servers[i].Manager.Dep.Id = a.DepId
		servers[i].Manager.Dep.Name = a.DepName
	}
	count, _ := orm.NewOrm().QueryTable("server").Count()
	return count, servers, true
}

//服务器检查
func ServerCheck(maId int, hardPosition string) (int64, bool) {
	o := orm.NewOrm().QueryTable("server")
	if maId != 0 {
		o = o.Filter("manager_id", maId)
	}
	if hardPosition != "" {
		o = o.Filter("hard_position", hardPosition)
	}
	count, err := o.Count()
	if err != nil {
		beego.Error("server check error : ", err)
		return 0, false
	}
	return count, true
}
