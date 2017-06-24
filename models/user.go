package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
)

type User struct {
	Id         int
	UserName   string    `orm:"unique;size(20)"`             // 用户名
	Password   string    `orm:"size(32)"`                    // 密码
	Salt       string    `orm:"size(10)"`                    // 密码盐
	Sex        int       `orm:"default(0)"`                  // 性别
	Email      string    `orm:"size(50)"`                    // 邮箱
	LastLogin  time.Time `orm:"null;type(datetime)"`         // 最后登录时间
	LastIp     string    `orm:"size(15)"`                    // 最后登录IP
	Status     int       `orm:"default(0)"`                  // 状态，0正常 -1禁用
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	RoleList   []Role    `orm:"-"`                           // 角色列表
}
type UserRole struct {
	UserId int // 用户id
	RoleId int // 角色id
}

func (this *User) table() string {
	return tableName("user")
}
func (m *User) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func (this *User) GetUser(userId int, getRoleInfo bool) (*User, error) {
	user := &User{}
	user.Id = userId
	err := o.Read(user)
	if err == nil && getRoleInfo {
		user.RoleList, _ = this.GetUserRoleList(user.Id)
	}
	return user, err
}

func (this *User) GetUserRoleList(userId int) ([]Role, error) {
	var (
		roleRef  []*UserRole
		roleList []Role
		r        Role
	)
	sql := "SELECT role_id FROM " + tableName("user_role") + " WHERE user_id = ?"
	o.Raw(sql, userId).QueryRows(&roleRef)
	roleList = make([]Role, 0, len(roleRef))
	for _, v := range roleRef {

		role, err := r.GetRole(v.RoleId)
		if err == nil {
			roleList = append(roleList, *role)
		}
	}
	return roleList, nil
}

// 获取用户总数
func (this *User) GetTotal() (int64, error) {
	return o.QueryTable(this.table()).Count()
}

// 根据用户名获取用户信息
func (this *User) GetUserByName(userName string) (*User, error) {
	user := &User{}
	user.UserName = userName
	err := o.Read(user, "UserName")
	return user, err
}

// 更新用户信息
func (this *User) UpdateUser(user *User, fileds ...string) error {
	if len(fileds) < 1 {
		return errors.New("更新字段不能为空")
	}
	_, err := o.Update(user, fileds...)
	return err
}

// 修改密码
func (this *User) ModifyPassword(userId int, password string) error {
	user, err := this.GetUser(userId, false)
	if err != nil {
		return err
	}
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = Md5([]byte(password + user.Salt))
	_, err = o.Update(user, "Salt", "Password")
	return err
}
