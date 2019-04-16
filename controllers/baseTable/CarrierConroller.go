package baseTable

import (
	"CcsuWebMonitor/controllers"
	"CcsuWebMonitor/models/baseClass"
	"encoding/json"

	"github.com/astaxie/beego"
)

type CarrierController struct {
	controllers.BaseController
	carrier baseClass.CarrierList
}

type carrierSearch struct {
	Carrier string
	Address string
}

//@router /ccsu/base_table/carrier_get [*]
func (c *CarrierController) GetCarrier() {
	if controllers.LoginCheck(&c.BaseController) {
		var search carrierSearch
		json.Unmarshal(c.Ctx.Input.RequestBody, &search)
		if count, carriers, err := baseClass.CarrierGet(search.Carrier, search.Address); err != nil {
			beego.Error("Get carriers error : ", err)
			c.Data["json"] = controllers.ResultJson(7001, "获取失败")
			c.ServeJSON()
			return
		} else {
			result := controllers.ResultJson(6000, "")
			result["count"] = count
			result["carriers"] = carriers
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}
}

//@router /ccsu/base_table/carrier_add [*]
func (c *CarrierController) AddCarrier() {
	if controllers.LoginCheck(&c.BaseController) {
		json.Unmarshal(c.Ctx.Input.RequestBody, &c.carrier)
		err := c.carrier.CarrierAdd()
		if err != nil {
			beego.Error("Add carrier error : ", err)
			c.Data["json"] = controllers.ResultJson(7001, "该IP已存在")
		} else {
			c.Data["json"] = controllers.ResultJson(7000, "添加成功")
		}
		c.ServeJSON()
		return
	}
}

//@router /ccsu/base_table/carrier_update [*]
func (c *CarrierController) UpDateCarrier() {
	if controllers.LoginCheck(&c.BaseController) { //确认是否登陆
		json.Unmarshal(c.Ctx.Input.RequestBody, &c.carrier)
		flag, err := c.carrier.CarrierUpdate()
		switch flag {
		case 0:
			beego.Error("Update carrier error : ", err)
			c.Data["json"] = controllers.ResultJson(7001, "修改失败")
			break
		case 1:
			c.Data["json"] = controllers.ResultJson(7001, "IP已重复")
			break
		case 2:
			c.Data["json"] = controllers.ResultJson(7000, "修改成功")
			break
		}
		c.ServeJSON()
		return
	}
}

//@router /ccsu/base_table/carrier_delete [*]
func (c *CarrierController) DeleteCarrier() {
	if controllers.LoginCheck(&c.BaseController) {
		//department := s.departmentGet()
		json.Unmarshal(c.Ctx.Input.RequestBody, &c.carrier)
		if err := baseClass.CarrierDelete(c.carrier.Id); err != nil {
			c.Data["json"] = controllers.ResultJson(7001, "删除失败")
		} else {
			c.Data["json"] = controllers.ResultJson(7000, "删除成功")
		}
		c.ServeJSON()
	}
}
