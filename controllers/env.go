package controllers

import (
	"gitdeploy/models"
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
