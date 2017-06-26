package controllers

import (
	"gitdeploy/models"

	"github.com/astaxie/beego"
)

type MailTplController struct {
	BaseController
}

// 模板列表
func (this *MailTplController) List() {
	var (
		mails models.MailTpl
	)
	list, _ := mails.GetMailTplList()
	this.Data["pageTitle"] = "邮件模板"
	this.Data["list"] = list
	this.display()
}

// 添加模板
func (this *MailTplController) Add() {
	var (
		mails models.MailTpl
	)
	if this.isPost() {
		name := this.GetString("name")
		subject := this.GetString("subject")
		content := this.GetString("content")
		mailTo := this.GetString("mail_to")
		mailCc := this.GetString("mail_cc")

		if name == "" || subject == "" || content == "" {
			this.showMsg("模板名称、邮件主题、邮件内容不能为空", MSG_ERR)
		}

		tpl := new(models.MailTpl)
		tpl.UserId = this.auth.GetUserId()
		tpl.Name = name
		tpl.Subject = subject
		tpl.Content = content
		tpl.MailTo = mailTo
		tpl.MailCc = mailCc
		err := mails.AddMailTpl(tpl)
		this.checkError(err)

		this.redirect(beego.URLFor("MailTplController.List"))
	}

	this.Data["pageTitle"] = "添加模板"
	this.display()
}

// 编辑模板
func (this *MailTplController) Edit() {
	var (
		mails models.MailTpl
	)
	id, _ := this.GetInt("id")
	tpl, err := mails.GetMailTpl(id)
	this.checkError(err)

	if this.isPost() {
		name := this.GetString("name")
		subject := this.GetString("subject")
		content := this.GetString("content")
		mailTo := this.GetString("mail_to")
		mailCc := this.GetString("mail_cc")
		if name == "" || subject == "" || content == "" {
			this.showMsg("模板名称、邮件主题、邮件内容不能为空", MSG_ERR)
		}

		tpl.Name = name
		tpl.Subject = subject
		tpl.Content = content
		tpl.MailTo = mailTo
		tpl.MailCc = mailCc
		err := mails.SaveMailTpl(tpl)
		this.checkError(err)

		this.redirect(beego.URLFor("MailTplController.List"))
	}

	this.Data["pageTitle"] = "修改模板"
	this.Data["tpl"] = tpl
	this.display()
}

// 删除模板
func (this *MailTplController) Del() {
	var (
		mails models.MailTpl
	)
	id, _ := this.GetInt("id")

	err := mails.DelMailTpl(id)
	this.checkError(err)

	this.redirect(beego.URLFor("MailTplController.List"))
}
