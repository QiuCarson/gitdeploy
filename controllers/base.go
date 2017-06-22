package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BaseController struct {
	beego.Controller
}

var (
	o orm.Ormer
)

func (this *BaseController) Prepare() {
	this.initAuth()
}

func (this *BaseController) initAuth() {
	token := this.Ctx.GetCookie("auth")
}
