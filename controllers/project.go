package controllers

import (
	"fmt"
	"gitdeploy/libs"
	"gitdeploy/models"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/lisijie/gopub/app/service"
)

type ProjectController struct {
	BaseController
}

// 项目列表
func (this *ProjectController) List() {
	var (
		projects models.Project
	)
	page, _ := strconv.Atoi(this.GetString("page"))
	if page < 1 {
		page = 1
	}

	count, _ := projects.GetTotal()
	list, _ := projects.GetList(page, this.pageSize)

	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("ProjectController.List"), true).ToString()
	this.Data["pageTitle"] = "项目列表"
	this.display()
}

// 添加项目
func (this *ProjectController) Add() {
	var (
		projects models.Project
		actions  models.Action
		servers  models.Server
		auths    models.Auth
	)
	if this.isPost() {
		p := &models.Project{}
		p.Name = this.GetString("project_name")
		p.Domain = this.GetString("project_domain")
		p.RepoUrl = this.GetString("repo_url")
		p.AgentId, _ = this.GetInt("agent_id")
		p.IgnoreList = this.GetString("ignore_list")
		p.BeforeShell = this.GetString("before_shell")
		p.AfterShell = this.GetString("after_shell")
		p.TaskReview, _ = this.GetInt("task_review")
		if v, _ := this.GetInt("create_verfile"); v > 0 {
			p.CreateVerfile = 1
		} else {
			p.CreateVerfile = 0
		}
		p.VerfilePath = strings.Replace(this.GetString("verfile_path"), ".", "", -1)

		if err := this.validProject(p); err != nil {
			this.showMsg(err.Error(), MSG_ERR)
		}

		err := projects.AddProject(p)
		this.checkError(err)

		// 克隆仓库
		go projects.CloneRepo(p.Id)

		actions.Add("add_project", auths.GetUserName(), "project", p.Id, "")

		this.redirect(beego.URLFor("ProjectController.List"))
	}

	agentList, err := servers.GetAgentList(1, -1)
	this.checkError(err)
	this.Data["pageTitle"] = "添加项目"
	this.Data["agentList"] = agentList
	this.display()
}

// 验证提交
func (this *ProjectController) validProject(p *models.Project) error {
	var (
		servers models.Server
	)
	errorMsg := ""
	if p.Name == "" {
		errorMsg = "请输入项目名称"
	} else if p.Domain == "" {
		errorMsg = "请输入项目标识"
	} else if p.RepoUrl == "" {
		errorMsg = "请输入仓库地址"
	} else if p.AgentId == 0 {
		errorMsg = "请选择跳板机"
	} else {
		agent, err := servers.GetServer(p.AgentId)
		if err != nil {
			return err
		}
		addr := fmt.Sprintf("%s:%d", agent.Ip, agent.SshPort)
		serv := libs.NewServerConn(addr, agent.SshUser, agent.SshKey)
		workPath := fmt.Sprintf("%s/%s", agent.WorkDir, p.Domain)

		if err := serv.TryConnect(); err != nil {
			errorMsg = "无法连接到跳板机: " + err.Error()
		} else if _, err := serv.RunCmd("mkdir -p " + workPath); err != nil {
			errorMsg = "无法创建跳板机工作目录: " + err.Error()
		}
		serv.Close()
	}

	if errorMsg != "" {
		return fmt.Errorf(errorMsg)
	}
	return nil
}

// 删除项目
func (this *ProjectController) Del() {
	var (
		projects models.Project
	)

	id, _ := this.GetInt("id")

	err := projects.DeleteProject(id)
	this.checkError(err)

	service.ActionService.Add("del_project", this.auth.GetUserName(), "project", id, "")

	this.redirect(beego.URLFor("ProjectController.List"))
}
