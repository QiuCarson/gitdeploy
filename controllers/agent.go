package controllers

import (
	"errors"
	"fmt"
	"gitdeploy/libs"
	"gitdeploy/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type AgentController struct {
	BaseController
}

// 列表
func (this *AgentController) List() {
	var (
		servers models.Server
	)
	page, _ := strconv.Atoi(this.GetString("page"))
	if page < 1 {
		page = 1
	}
	count, err := servers.GetTotal(models.SERVER_TYPE_AGENT)
	this.checkError(err)
	serverList, err := servers.GetAgentList(page, this.pageSize)
	this.checkError(err)

	this.Data["count"] = count
	this.Data["list"] = serverList
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AgentController.List"), true).ToString()
	this.Data["pageTitle"] = "跳板机列表"
	this.display()
}

// 添加
func (this *AgentController) Add() {
	if this.isPost() {
		server := &models.Server{}
		server.TypeId = models.SERVER_TYPE_AGENT
		server.Ip = this.GetString("server_ip")
		server.Area = this.GetString("area")
		server.SshPort, _ = this.GetInt("ssh_port")
		server.SshUser = this.GetString("ssh_user")
		server.SshPwd = this.GetString("ssh_pwd")
		server.SshKey = this.GetString("ssh_key")
		server.WorkDir = this.GetString("work_dir")
		server.Description = this.GetString("description")
		err := this.validServer(server)
		this.checkError(err)
		err = server.AddServer(server)
		this.checkError(err)
		//service.ActionService.Add("add_agent", this.auth.GetUserName(), "server", server.Id, server.Ip)
		this.redirect(beego.URLFor("AgentController.List"))
	}

	this.Data["pageTitle"] = "添加跳板机"
	this.display()
}

func (this *AgentController) validServer(server *models.Server) error {
	valid := validation.Validation{}
	valid.Required(server.Ip, "ip").Message("请输入服务器IP")
	valid.Range(server.SshPort, 1, 65535, "ssh_port").Message("SSH端口无效")
	valid.Required(server.SshUser, "ssh_user").Message("SSH用户名不能为空")
	valid.Required(server.WorkDir, "work_dir").Message("工作目录不能为空")
	valid.Required(server.SshPwd, "ssh_pwd").Message("SSH密码不能为空")
	valid.IP(server.Ip, "ip").Message("服务器IP无效")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	if server.SshKey != "" && !libs.IsFile(libs.RealPath(server.SshKey)) {
		return errors.New("SSH Key不存在:" + server.SshKey)
	}

	addr := fmt.Sprintf("%s:%d", server.Ip, server.SshPort)
	serv := libs.NewServerConn(addr, server.SshUser, server.SshKey, server.SshPwd)

	if err := serv.TryConnect(); err != nil {
		return errors.New("无法连接到跳板机: " + err.Error())
	} else if _, err := serv.RunCmd("mkdir -p " + server.WorkDir); err != nil {
		return errors.New("无法创建跳板机工作目录: " + err.Error())
	}
	serv.Close()

	return nil
}

// 编辑
func (this *AgentController) Edit() {
	id, _ := this.GetInt("id")
	server, err := new(models.Server).GetServer(id, models.SERVER_TYPE_AGENT)
	this.checkError(err)

	if this.isPost() {
		server.Ip = this.GetString("server_ip")
		server.Area = this.GetString("area")
		server.SshPort, _ = this.GetInt("ssh_port")
		server.SshUser = this.GetString("ssh_user")
		server.SshPwd = this.GetString("ssh_pwd")
		server.SshKey = this.GetString("ssh_key")
		server.WorkDir = this.GetString("work_dir")
		server.Description = this.GetString("description")
		err := this.validServer(server)
		this.checkError(err)
		err = new(models.Server).UpdateServer(server)
		this.checkError(err)
		//service.ActionService.Add("edit_agent", this.auth.GetUserName(), "server", server.Id, server.Ip)
		this.redirect(beego.URLFor("AgentController.List"))
	}

	this.Data["pageTitle"] = "编辑跳板机"
	this.Data["server"] = server
	this.display()
}

// 删除
func (this *AgentController) Del() {
	id, _ := this.GetInt("id")

	_, err := new(models.Server).GetServer(id, models.SERVER_TYPE_AGENT)
	this.checkError(err)

	err = new(models.Server).DeleteServer(id)
	this.checkError(err)
	//service.ActionService.Add("del_agent", this.auth.GetUserName(), "server", id, "")

	this.redirect(beego.URLFor("AgentController.List"))
}

// 项目列表
func (this *AgentController) Projects() {
	id, _ := this.GetInt("id")
	server, err := new(models.Server).GetServer(id, models.SERVER_TYPE_AGENT)
	this.checkError(err)
	envList, err := new(models.EnvServer).GetEnvListByServerId(id)
	this.checkError(err)

	result := make(map[int]map[string]interface{})
	for _, env := range envList {
		if _, ok := result[env.ProjectId]; !ok {
			project, err := new(models.Project).GetProject(env.ProjectId)
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
