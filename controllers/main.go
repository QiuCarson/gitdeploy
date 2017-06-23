package controllers

import (
	"fmt"
	"gitdeploy/models"
	"os"
	"runtime"
	"time"

	"github.com/astaxie/beego"
)

type MainController struct {
	BaseController
}

// 首页
func (this *MainController) Index() {
	var (
		tasks    models.Task
		servers  models.Server
		users    models.User
		projects models.Project
		actions  models.Action
	)
	this.Data["pageTitle"] = "系统概况"
	projectsStat := tasks.GetProjectPubStat()

	popProjects := make([]map[string]interface{}, 0, 4)
	for k, stat := range projectsStat {
		projectInfo, err := projects.GetProject(stat["project_id"])
		if err != nil {
			continue
		}
		if k > 4 {
			break
		}
		info := make(map[string]interface{})
		info["project_name"] = projectInfo.Name
		info["version"] = projectInfo.Version
		info["version_time"] = beego.Date(projectInfo.VersionTime, "Y-m-d H:i:s")
		info["count"] = stat["count"]
		popProjects = append(popProjects, info)
	}

	feeds, _ := actions.GetList(1, 7)
	this.Data["feeds"] = feeds
	this.Data["serverNum"], _ = servers.GetTotal(models.SERVER_TYPE_NORMAL)
	this.Data["userNum"], _ = users.GetTotal()
	this.Data["projectNum"], _ = projects.GetTotal()
	this.Data["pubNum"], _ = tasks.GetPubTotal()
	this.Data["popProjects"] = popProjects
	this.Data["hostname"], _ = os.Hostname()
	this.Data["gover"] = runtime.Version()
	this.Data["os"] = runtime.GOOS
	this.Data["goroutineNum"] = runtime.NumGoroutine()
	this.Data["cpuNum"] = runtime.NumCPU()
	this.Data["arch"] = runtime.GOARCH
	this.Data["dbVerson"] = models.DBVersion()
	this.Data["dataDir"] = beego.AppConfig.String("data_dir")

	up, day, hour, min, sec := this.getUptime()
	this.Data["uptime"] = fmt.Sprintf("%s，已运行 %d天 %d小时 %d分钟 %d秒", beego.Date(up, "Y-m-d H:i:s"), day, hour, min, sec)

	this.display()
}

func (this *MainController) getUptime() (up time.Time, day, hour, min, sec int) {
	ts, _ := beego.AppConfig.Int64("up_time")
	up = time.Unix(ts, 0)
	uptime := int(time.Now().Sub(up) / time.Second)
	if uptime >= 86400 {
		day = uptime / 86400
		uptime %= 86400
	}
	if uptime >= 3600 {
		hour = uptime / 3600
		uptime %= 3600
	}
	if uptime >= 60 {
		min = uptime / 60
		uptime %= 60
	}
	sec = uptime
	return
}

// 登录
func (this *MainController) Login() {
	var (
	//actions models.Action
	//auths   models.Auth
	)
	if this.userId > 0 {
		this.redirect("/")
	}
	/*beego.ReadFromRequest(&this.Controller)
	if this.isPost() {
		flash := beego.NewFlash()
		username := this.GetString("username")
		password := this.GetString("password")
		remember := this.GetString("remember")
		if username != "" && password != "" {
			token, err := auths.Login(username, password)
			if err != nil {
				flash.Error(err.Error())
				flash.Store(&this.Controller)
				this.redirect("/login")
			} else {
				if remember == "yes" {
					this.Ctx.SetCookie("auth", token, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", token)
				}
				actions.Login(username, auths.GetUserId(), this.getClientIp())
				this.redirect(beego.URLFor(".Index"))
			}

		}
	}*/

	this.TplName = "main/login.html"
}

// 退出登录
func (this *MainController) Logout() {
	var (
		actions models.Action
		auths   models.Auth
	)
	actions.Logout(auths.GetUser().UserName, auths.GetUserId(), this.getClientIp())
	auths.Logout()
	this.Ctx.SetCookie("auth", "")
	this.redirect(beego.URLFor(".Login"))
}

// 个人信息
func (this *MainController) Profile() {

}
