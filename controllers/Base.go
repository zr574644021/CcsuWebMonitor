package controllers

import (
	"CcsuWebMonitor/models/baseClass"
	"CcsuWebMonitor/models/responseClass"
	"fmt"
	"net"
	"net/http"
	"time"
)

var Client = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*3)
			if err != nil {
				fmt.Println("dail timeout", err)
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 2,
	},
}

/*func init() {
	Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 2,
		},
	}
}*/

func ResultJson(code int, message string) (result map[string]interface{}) {
	result = make(map[string]interface{})
	result["status"] = code
	if message != "" {
		result["message"] = message
	}
	return result
}

func LoginCheck(c *BaseController) bool {
	username := c.GetSession("username")
	password := c.GetSession("password")    //密文形式的密码
	if username == nil || password == nil { //尝试获取session,获取不到则返回未登录
		c.Data["json"] = ResultJson(7002, "未登录")
		c.ServeJSON()
		return false
	} else {
		var user baseClass.User
		user.UserName = username.(string) //GetSession返回的为interface{} 需要使用.(string)强制转换成string类型
		user.Password = password.(string)
		if user.Login_matching_crypte() { //验证session中的用户名和密码是否正确
			return true
		} else {
			c.Data["json"] = ResultJson(7003, "session过期") //错误返回session异常
			c.ServeJSON()
		}
	}
	return false
}

func HttpGetTime(url string, get chan responseClass.StatusTime) {
	//将计算好的时间差和状态码放入通道
	var statustime responseClass.StatusTime
	time1 := time.Now() //获取请求前的时间
	res, err := Client.Get("http://" + url)
	time2 := time.Now() //获取请求后的时间
	if err != nil {
		//if get request error.content of error input in status
		statustime.Status = fmt.Sprintf("%s", err)
		get <- statustime
		close(get)
		return
	}

	statustime.Time = time2.Sub(time1).Seconds()
	statustime.Status = res.Status
	get <- statustime
	close(get)
}
