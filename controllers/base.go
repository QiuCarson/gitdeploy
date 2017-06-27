package controllers

import (
	"encoding/json"
	"gitdeploy/models"
	"io/ioutil"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
)

const (
	MSG_OK       = 0  // ajax输出错误码，成功
	MSG_ERR      = -1 // 错误
	MSG_REDIRECT = -2 // 重定向
)

type BaseController struct {
	beego.Controller
	auth           *models.Auth
	userId         int    // 当前登录的用户id
	controllerName string // 控制器名
	actionName     string // 动作名
	pageSize       int    // 默认分页大小
	lang           string // 当前语言环境
}

var (
	o orm.Ormer
)

func (this *BaseController) Prepare() {
	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	this.actionName = strings.ToLower(actionName)
	this.initAuth()
	this.initLang()
}

func (this *BaseController) initAuth() {
	//var p *models.Auth
	token := this.Ctx.GetCookie("auth")
	this.auth = new(models.Auth)

	this.auth.Init(token)

	this.userId = this.auth.GetUserId()

	if !this.auth.IsLogined() {
		if this.controllerName != "main" ||
			(this.controllerName == "main" && this.actionName != "logout" && this.actionName != "login") {
			this.redirect(beego.URLFor("MainController.Login"))
		}
	} else {
		if !this.auth.HasAccessPerm(this.controllerName, this.actionName) {
			this.showMsg("您没有执行该操作的权限", MSG_ERR)
		}
	}

}

//国际化
func (this *BaseController) initLang() {
	this.lang = "zh-CN"
	this.Data["lang"] = this.lang
	if !i18n.IsExist(this.lang) {
		if err := i18n.SetMessage(this.lang, beego.AppPath+"/conf/locale_"+this.lang+".ini"); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}

	}
}

// 重定向
func (this *BaseController) redirect(url string) {
	if this.IsAjax() {
		this.showMsg("", MSG_REDIRECT, url)
	} else {
		this.Redirect(url, 302)
		this.StopRun()
	}
}

// 提示消息
func (this *BaseController) showMsg(msg string, msgno int, redirect ...string) {
	out := make(map[string]interface{})
	out["status"] = msgno
	out["msg"] = msg
	out["redirect"] = ""
	if len(redirect) > 0 {
		out["redirect"] = redirect[0]
	}

	if this.IsAjax() {
		this.jsonResult(out)
	} else {
		for k, v := range out {
			this.Data[k] = v
		}
		this.display("error/message")
		this.Render()
		this.StopRun()
	}
}

// 输出json
func (this *BaseController) jsonResult(out interface{}) {
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

//渲染模版
func (this *BaseController) display(tpl ...string) {

	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = this.controllerName + "/" + this.actionName + ".html"
	}

	this.Layout = "layout/layout.html"
	this.TplName = tplname

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "layout/sections/header.html"
	this.LayoutSections["Footer"] = "layout/sections/footer.html"
	this.LayoutSections["Navbar"] = "layout/sections/navbar.html"
	this.LayoutSections["Sidebar"] = "layout/sections/sidebar.html"

	user := this.auth.GetUser()

	this.Data["version"] = beego.AppConfig.String("version")
	this.Data["curRoute"] = this.controllerName + "." + this.actionName
	this.Data["loginUserId"] = user.Id
	this.Data["loginUserName"] = user.UserName
	this.Data["loginUserSex"] = user.Sex
	this.Data["menuList"] = this.getMenuList()
}

// 功能菜单
func (this *BaseController) getMenuList() []Menu {
	var menuList []Menu
	allMenu := make([]Menu, 0)
	content, err := ioutil.ReadFile("conf/menu.json")
	if err == nil {
		err := json.Unmarshal(content, &allMenu)
		if err != nil {
			beego.Error(err.Error())
		}
	}
	menuList = make([]Menu, 0)
	for _, menu := range allMenu {
		subs := make([]SubMenu, 0)
		for _, sub := range menu.Submenu {
			route := strings.Split(sub.Route, ".")
			if this.auth.HasAccessPerm(route[0], route[1]) {
				subs = append(subs, sub)
			}
		}
		if len(subs) > 0 {
			menu.Submenu = subs
			menuList = append(menuList, menu)
		}
	}
	//menuList = allMenu

	return menuList
}

// 是否POST提交
func (this *BaseController) isPost() bool {
	return this.Ctx.Request.Method == "POST"
}

//获取用户IP地址
func (this *BaseController) getClientIp() string {
	if p := this.Ctx.Input.Proxy(); len(p) > 0 {
		return p[0]
	}
	return this.Ctx.Input.IP()
}

// 错误检查
func (this *BaseController) checkError(err error) {
	if err != nil {
		this.showMsg(err.Error(), MSG_ERR)
	}
}
