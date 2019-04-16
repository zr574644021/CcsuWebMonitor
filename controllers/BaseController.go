package controllers

import (
	"encoding/json"

	"CcsuWebMonitor/models/baseClass"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	user baseClass.User
}

//@router /ccsu/login [*]
func (c *BaseController) LoginPost() {
	json.Unmarshal(c.Ctx.Input.RequestBody, &c.user) //json格式解析
	switch c.user.Login_matching() {
	case 0:
		c.SetSession("username", c.user.UserName)
		c.SetSession("password", baseClass.Encrypt(c.user.Password))
		result := ResultJson(7000, "")
		result["user_name"] = c.user.UserName
		c.Data["json"] = result                           //返回登陆成功
		c.user.Login_record("", c.Ctx.Request.RemoteAddr) //写入登陆记录
	case 1:
		c.Data["json"] = ResultJson(7001, "账号不存在") //账号或密码错误
		//c.user.Login_record(c.user.Password, c.Ctx.Request.RemoteAddr) //将错误的密码记录
	case 2:
		c.Data["json"] = ResultJson(7002, "密码错误")                      //账号或密码错误
		c.user.Login_record(c.user.Password, c.Ctx.Request.RemoteAddr) //将错误的密码记录
	}
	c.ServeJSON()
}

//@router /ccsu/logout [*]
func (c *BaseController) LogoutPost() {
	username := c.GetSession("username")
	password := c.GetSession("password")    //密文形式的密码
	if username == nil || password == nil { //尝试获取session,获取不到则返回未登录
		c.Data["json"] = ResultJson(7002, "未登录")
		c.ServeJSON()
		return
	} else {
		c.DelSession("username")
		c.DelSession("password")
		c.Data["json"] = ResultJson(7000, "")
		c.ServeJSON()
		return
	}
}
