package controllers

import (
	"encoding/json"
	"gitdeploy/models"
	"strconv"

	"github.com/astaxie/beego"
)

type EnvController struct {
	BaseController
}

func (this *EnvController) List() {
	var (
		envs models.Env
	)
	projectId, _ := this.GetInt("project_id")
	envList, _ := envs.GetEnvListByProjectId(projectId)
	this.Data["pageTitle"] = "发布环境配置"
	this.Data["projectId"] = projectId
	this.Data["envList"] = envList
	this.display()
}
func (this *EnvController) Add() {
	var (
		projects   models.Project
		servers    models.Server
		envservers models.EnvServer
		mailTpls   models.MailTpl
	)
	projectId, _ := this.GetInt("project_id")

	project, err := projects.GetProject(projectId)
	this.checkError(err)

	if this.isPost() {
		env := new(models.Env)
		env.ProjectId = project.Id
		env.Name = this.GetString("name")
		env.SshUser = this.GetString("ssh_user")
		env.SshPort = this.GetString("ssh_port")
		env.SshKey = this.GetString("ssh_key")
		env.PubDir = this.GetString("pub_dir")
		env.BeforeShell = this.GetString("before_shell")
		env.AfterShell = this.GetString("after_shell")
		env.SendMail, _ = this.GetInt("send_mail")
		env.MailTplId, _ = this.GetInt("mail_tpl_id")
		env.MailTo = this.GetString("mail_to")
		env.MailCc = this.GetString("mail_cc")

		if env.Name == "" || env.SshUser == "" || env.SshPort == "" || env.SshKey == "" || env.PubDir == "" {
			this.showMsg("环境名称、SSH帐号、SSH端口、SSH KEY路径、发布目录不能为空。", MSG_ERR)
		}

		serverIds := this.GetStrings("serverIds")
		if len(serverIds) < 1 {
			this.showMsg("请选择服务器", MSG_ERR)
		}

		if env.SendMail > 0 {
			if env.MailTplId == 0 {
				this.showMsg("请选择邮件模板", MSG_ERR)
			}
		}

		env.ServerList = make([]models.Server, 0, len(serverIds))
		for _, v := range serverIds {
			if sid, _ := strconv.Atoi(v); sid > 0 {
				if sv, err := servers.GetServer(sid); err == nil {
					env.ServerList = append(env.ServerList, *sv)
				} else {
					this.showMsg("服务器ID不存在: "+v, MSG_ERR)
				}
			}
		}
		if err := envservers.AddEnv(env); err != nil {
			this.checkError(err)
		}

		this.redirect(beego.URLFor("EnvController.List", "project_id", projectId))
	}

	this.Data["serverList"], _ = servers.GetServerList(1, -1)
	this.Data["mailTplList"], _ = mailTpls.GetMailTplList()
	this.Data["project"] = project
	this.Data["pageTitle"] = "添加发布环境"
	this.display()
}

func (this *EnvController) Edit() {
	var (
		servers    models.Server
		envservers models.EnvServer
		mailTpls   models.MailTpl
	)
	id, _ := this.GetInt("id")

	env, err := envservers.GetEnv(id)
	this.checkError(err)

	if this.isPost() {
		env.Name = this.GetString("name")
		env.SshUser = this.GetString("ssh_user")
		env.SshPort = this.GetString("ssh_port")
		env.SshKey = this.GetString("ssh_key")
		env.PubDir = this.GetString("pub_dir")
		env.BeforeShell = this.GetString("before_shell")
		env.AfterShell = this.GetString("after_shell")
		env.SendMail, _ = this.GetInt("send_mail")
		env.MailTplId, _ = this.GetInt("mail_tpl_id")
		env.MailTo = this.GetString("mail_to")
		env.MailCc = this.GetString("mail_cc")

		if env.Name == "" || env.SshUser == "" || env.SshPort == "" || env.SshKey == "" || env.PubDir == "" {
			this.showMsg("环境名称、SSH帐号、SSH端口、SSH KEY路径、发布目录不能为空。", MSG_ERR)
		}

		serverIds := this.GetStrings("serverIds")
		if len(serverIds) < 1 {
			this.showMsg("请选择服务器", MSG_ERR)
		}

		if env.SendMail > 0 {
			if env.MailTplId == 0 {
				this.showMsg("请选择邮件模板", MSG_ERR)
			}
		}

		env.ServerList = make([]models.Server, 0, len(serverIds))
		for _, v := range serverIds {
			if sid, _ := strconv.Atoi(v); sid > 0 {
				if sv, err := servers.GetServer(sid); err == nil {
					env.ServerList = append(env.ServerList, *sv)
				} else {
					this.showMsg("服务器ID不存在: "+v, MSG_ERR)
				}
			}
		}

		envservers.SaveEnv(env)

		this.redirect(beego.URLFor("EnvController.List", "project_id", env.ProjectId))
	}

	serverList, _ := servers.GetServerList(1, -1)

	serverIds := make([]int, 0, len(env.ServerList))
	for _, v := range env.ServerList {
		serverIds = append(serverIds, v.Id)
	}

	jsonData, err := json.Marshal(serverIds)
	this.checkError(err)
	mailTplList, _ := mailTpls.GetMailTplList()

	this.Data["serverList"] = serverList
	this.Data["mailTplList"] = mailTplList
	this.Data["serverIds"] = string(jsonData)
	this.Data["env"] = env
	this.Data["pageTitle"] = "编辑发布环境"
	this.display()
}

func (this *EnvController) Del() {
	var (
		envservers models.EnvServer
	)
	id, _ := this.GetInt("id")
	envservers.DeleteEnv(id)
	this.redirect(beego.URLFor("EnvController.List", "project_id", id))
}
