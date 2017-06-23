package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Auth struct {
	loginUser *User
	openPerm  map[string]bool // 公开的权限
	permMap   map[string]bool // 当前用户权限表
}

// 初始化开放权限
func (this *Auth) initOpenPerm() {
	this.openPerm = map[string]bool{
		"main.index":        true,
		"main.profile":      true,
		"main.login":        true,
		"main.logout":       true,
		"main.getpubstat":   true,
		"project.clone":     true,
		"project.getstatus": true,
		"task.gettags":      true,
		"task.getstatus":    true,
		"task.startpub":     true,
	}

}

func (this *Auth) Init(token string) {

	var u User

	this.initOpenPerm()

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

// 获取当前登录的用户id
func (this *Auth) GetUserId() int {
	if this.IsLogined() {
		return this.loginUser.Id

	}
	return 0
}

// 获取当前登录的用户对象
func (this *Auth) GetUser() *User {
	return this.loginUser
}

// 检查是否有某个权限
func (this *Auth) HasAccessPerm(module, action string) bool {
	key := module + "." + action
	if !this.IsLogined() {
		return false
	}
	if this.loginUser.Id == 1 || this.isOpenPerm(key) {
		return true
	}
	if _, ok := this.permMap[key]; ok {
		return true
	}
	return false
}

// 是否公开访问的操作
func (this *Auth) isOpenPerm(key string) bool {
	if _, ok := this.openPerm[key]; ok {
		return true
	}
	return false
}

// 用户登录
func (this *Auth) Login(userName, password string) (string, error) {
	var users User

	user, err := users.GetUserByName(userName)
	if err != nil {
		if err == orm.ErrNoRows {
			return "", errors.New("帐号或密码错误")
		} else {
			return "", errors.New("系统错误")
		}
	}

	if user.Password != Md5([]byte(password+user.Salt)) {
		return "", errors.New("帐号或密码错误")
	}
	if user.Status == -1 {
		return "", errors.New("该帐号已禁用")
	}

	user.LastLogin = time.Now()
	users.UpdateUser(user, "LastLogin")
	this.loginUser = user

	token := fmt.Sprintf("%d|%s", user.Id, Md5([]byte(user.Password+user.Salt)))
	return token, nil
}

// 退出登录
func (this *Auth) Logout() error {
	return nil
}
