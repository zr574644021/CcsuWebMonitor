package webMonitor

import (
	"CcsuWebMonitor/controllers"
	"CcsuWebMonitor/models/responseClass"

	"encoding/json"

	"time"

	"CcsuWebMonitor/models/baseClass"
)

type UrlRecord struct {
	Url          string
	Responsetime float64
	Time
}

type Time struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
}

type UrlNow struct {
	Id           int
	Url          string
	WebName      string
	SleepTime    float64
	ResponseTime float64
	Status       string
	ResCode      int
}

type urllist struct {
	Id        int
	Url       string
	WebName   string
	SleepTime float64
}

type HttpController struct {
	controllers.BaseController
	baseClass.WebSiteList
}

//@router /ccsu/http_monitor/url_monitor [post]
func (s *HttpController) HttpMonitorQuery() {
	if controllers.LoginCheck(&s.BaseController) {
		statustime := make(chan responseClass.StatusTime)
		url := s.GetString("url")
		//open a go func to test outtime
		go controllers.HttpGetTime(url, statustime)
		select {
		//try to get data in chan
		case get := <-statustime:
			result := controllers.ResultJson(4003, "")
			result["time"] = get
			s.Data["json"] = result
			s.ServeJSON()
			return
			//get data in chan outtime
		case <-time.After(5 * time.Second):
			s.Data["json"] = controllers.ResultJson(-1, "")
			s.ServeJSON()
			return
		}
		close(statustime)
		s.Data["json"] = controllers.ResultJson(4001, "")
		s.ServeJSON()
		return
	}
}

/*//@router /ccsu/http_monitor/now [get]
func (s *HttpController) HttpMonitorNow() {
	if controllers.LoginCheck(&s.BaseController) {
		var record []UrlNow
		//get url all
		i, urllists, err := baseClass.WebSiteGet()
		if err != nil {
			s.Data["json"] = controllers.ResultJson(4004, "")
			s.ServeJSON()
			return
		}
		//make a chan get all returning data
		record = make([]UrlNow, i)
		var wg sync.WaitGroup
		for i, urllist := range urllists {
			wg.Add(1) //add a go func
			go func(i int, urlList baseClass.WebSiteList) {
				statustime := make(chan responseClass.StatusTime)
				go controllers.HttpGetTime(urlList.Url, statustime)
				var node UrlNow
				select {
				case get := <-statustime:
					node.ResponseTime = get.Time
					node.Status = get.Status
					node.ResCode = 200 //响应成功
				case <-time.After(2 * time.Second):
					node.ResponseTime = 0
					node.Status = "OutTime"
					node.ResCode = 400 //响应失败
				}
				node.Id = urlList.Id
				//node.SleepTime = urlList.SleepTime
				node.WebName = urlList.WebName
				node.Url = urlList.Url
				record[i] = node
				//close(statustime)
				wg.Done() // finish a go func
			}(i, urllist)
		}
		wg.Wait() // if have a go func not finish ,wait it
		result := controllers.ResultJson(4003, "")
		result["record"] = record
		s.Data["json"] = result
		s.ServeJSON()
		return
	}
}*/

//@router /ccsu/http_monitor/url_error_get [get]
func (s *HttpController) GetError() {
	if controllers.LoginCheck(&s.BaseController) {
		_, urlerror, err := responseClass.UrlErrorGet()
		if err != nil {
			s.Data["json"] = controllers.ResultJson(7001, "获取失败")
		} else {
			result := controllers.ResultJson(7000, "")
			result["urlerror"] = *urlerror
			s.Data["json"] = result
		}
		s.ServeJSON()
		return
	}
}

//@router /ccsu/http_monitor/url_error_delete [post]
func (s *HttpController) DeleteHttpError() {
	if controllers.LoginCheck(&s.BaseController) {
		var urldelete []string
		var urls []string
		json.Unmarshal(s.Ctx.Input.RequestBody, &urls)
		result := make(map[string]interface{})
		for _, url := range urls {
			if url == "" {
				s.Data["json"] = controllers.ResultJson(4001, "")
				s.ServeJSON()
				return
			}
			num, err := responseClass.UrlErrorDelete(url)
			if err != nil {
				result := controllers.ResultJson(4002, "")
				s.Data["json"] = result
				s.ServeJSON()
				return
			}
			urldelete = append(urldelete, url)
			result["num"] = result["num"].(int64) + num
			result["url_delete"] = urldelete
		}
		result = controllers.ResultJson(4003, "")
		s.Data["json"] = result
		s.ServeJSON()
		return
	}
}
