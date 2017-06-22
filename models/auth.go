package models

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

type Auth struct {
	loginUser *User
	permMap   map[string]bool // 当前用户权限表
}

func (this *Auth) Init(token string) {
	var u User
	arr := strings.Split(token, "|")
	beego.Trace("登录验证, token: ", token)
	if len(arr) == 2 {
		idstr, password := arr[0], arr[1]
		userId, _ := strconv.Atoi(idstr)
		if userId > 0 {
			user, err := u.GetUser(userId, true)
			if err == nil && password == Md5([]byte(user.Password+user.Salt)) {
				this.loginUser = user
				this.initPermMap()
				beego.Trace("验证成功，用户信息: ", user)
			}
		}
	}
}

// 初始化权限表
func (this *Auth) initPermMap() {
	this.permMap = make(map[string]bool)
	for _, role := range this.loginUser.RoleList {
		for _, perm := range role.PermList {
			this.permMap[perm.Key] = true
		}
	}
}

// 检查是否登录
func (this *Auth) IsLogined() bool {
	return this.loginUser != nil && this.loginUser.Id > 0
}
