package main

import (
	_ "gitdeploy/routers"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

func main() {
	beego.SetStaticPath("/assets", "assets")
	beego.AddFuncMap("i18n", i18n.Tr)

	beego.Run()
}
