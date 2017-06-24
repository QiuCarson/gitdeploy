package routers

import (
	"gitdeploy/controllers"

	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")

	beego.AutoRouter(&controllers.ProjectController{})

	beego.AutoRouter(&controllers.AgentController{})
	/*beego.AutoRouter(&controllers.TaskController{})
	beego.AutoRouter(&controllers.ServerController{})
	beego.AutoRouter(&controllers.EnvController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.RoleController{})
	beego.AutoRouter(&controllers.MailTplController{})
	beego.AutoRouter(&controllers.ReviewController{})
	beego.AutoRouter(&controllers.MainController{})*/
}
