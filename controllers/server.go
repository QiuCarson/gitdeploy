package controllers

import (
	"gitdeploy/libs"
	"gitdeploy/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type ServerController struct {
	BaseController
}

// 列表
func (this *ServerController) List() {
	var (
		servers models.Server
	)
	page, _ := strconv.Atoi(this.GetString("page"))
	if page < 1 {
		page = 1
	}
	count, err := servers.GetTotal(models.SERVER_TYPE_NORMAL)
	this.checkError(err)
	serverList, err := servers.GetServerList(page, this.pageSize)
	this.checkError(err)

	this.Data["count"] = count
	this.Data["list"] = serverList
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("ServerController.List"), true).ToString()
	this.Data["pageTitle"] = "服务器列表"
	this.display()
}

// 添加
func (this *ServerController) Add() {
	var (
		servers models.Server
	)
	if this.isPost() {
		valid := validation.Validation{}
		server := &models.Server{}
		server.TypeId = models.SERVER_TYPE_NORMAL
		server.Ip = this.GetString("server_ip")
		server.Area = this.GetString("area")
		server.Description = this.GetString("description")
		valid.Required(server.Ip, "ip").Message("请输入服务器IP")
		valid.IP(server.Ip, "ip").Message("服务器IP无效")
		if valid.HasErrors() {
			for _, err := range valid.Errors {
				this.showMsg(err.Message, MSG_ERR)
			}
		}

		if err := servers.AddServer(server); err != nil {
			this.showMsg(err.Error(), MSG_ERR)
		}

		this.redirect(beego.URLFor("ServerController.List"))
	}

	this.Data["pageTitle"] = "添加服务器"
	this.display()
}

// 编辑
func (this *ServerController) Edit() {
	var (
		servers models.Server
	)
	id, _ := this.GetInt("id")
	server, err := servers.GetServer(id, models.SERVER_TYPE_NORMAL)
	this.checkError(err)

	if this.isPost() {
		valid := validation.Validation{}
		ip := this.GetString("server_ip")
		server.Area = this.GetString("area")
		server.Description = this.GetString("description")
		valid.Required(ip, "ip").Message("请输入服务器IP")
		valid.IP(ip, "ip").Message("服务器IP无效")
		if valid.HasErrors() {
			for _, err := range valid.Errors {
				this.showMsg(err.Message, MSG_ERR)
			}
		}
		server.Ip = ip
		err := servers.UpdateServer(server)
		this.checkError(err)
		this.redirect(beego.URLFor("ServerController.List"))
	}

	this.Data["pageTitle"] = "编辑服务器"
	this.Data["server"] = server
	this.display()
}

// 删除
func (this *ServerController) Del() {
	var (
		servers models.Server
	)
	id, _ := this.GetInt("id")

	_, err := servers.GetServer(id, models.SERVER_TYPE_NORMAL)
	this.checkError(err)

	err = servers.DeleteServer(id)
	this.checkError(err)

	this.redirect(beego.URLFor("ServerController.List"))
}

// 项目列表
func (this *ServerController) Projects() {
	var (
		servers    models.Server
		envservers models.EnvServer
		projects   models.Project
	)
	id, _ := this.GetInt("id")
	server, err := servers.GetServer(id, models.SERVER_TYPE_NORMAL)
	this.checkError(err)
	envList, err := envservers.GetEnvListByServerId(id)
	this.checkError(err)

	result := make(map[int]map[string]interface{})
	for _, env := range envList {
		if _, ok := result[env.ProjectId]; !ok {
			project, err := projects.GetProject(env.ProjectId)
			if err != nil {
				continue
			}
			row := make(map[string]interface{})
			row["projectId"] = project.Id
			row["projectName"] = project.Name
			row["envName"] = env.Name
			result[env.ProjectId] = row
		} else {
			result[env.ProjectId]["envName"] = result[env.ProjectId]["envName"].(string) + ", " + env.Name
		}
	}

	this.Data["list"] = result
	this.Data["server"] = server
	this.Data["pageTitle"] = server.Ip + " 下的项目列表"
	this.display()
}
