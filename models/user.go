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

// 根据角色id获取用户列表
func (this *User) GetUserListByRoleId(roleId int) ([]User, error) {
	var users []User
	sql := "SELECT u.* FROM " + this.table() + " u JOIN " + tableName("user_role") + " r ON u.id = r.user_id WHERE r.role_id = ?"
	_, err := o.Raw(sql, roleId).QueryRows(&users)
	return users, err
}

// 设置用户角色
func (this *User) UpdateUserRoles(userId int, roleIds []int) error {
	if _, err := this.GetUser(userId, false); err != nil {
		return err
	}
	o.Raw("DELETE FROM "+tableName("user_role")+" WHERE user_id = ?", userId).Exec()
	for _, v := range roleIds {
		o.Raw("INSERT INTO "+tableName("user_role")+" (user_id, role_id) VALUES (?, ?)", userId, v).Exec()
	}
	return nil
}

// 分页获取用户列表
func (this *User) GetUserList(page, pageSize int, getRoleInfo bool) ([]User, error) {
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	var users []User
	qs := o.QueryTable(this.table())
	_, err := qs.OrderBy("id").Limit(pageSize, offset).All(&users)
	for k, user := range users {
		users[k].RoleList, _ = this.GetUserRoleList(user.Id)
	}

	return users, err
}

// 添加用户
func (this *User) AddUser(userName, email, password string, sex int) (*User, error) {
	if exists, _ := this.GetUserByName(userName); exists.Id > 0 {
		return nil, errors.New("用户名已存在")
	}

	user := &User{}
	user.UserName = userName
	user.Sex = sex
	user.Email = email
	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = Md5([]byte(password + user.Salt))
	// user.LastLogin = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	_, err := o.Insert(user)
	return user, err
}

// 删除用户
func (this *User) DeleteUser(userId int) error {
	if userId == 1 {
		return errors.New("不允许删除用户ID为1的用户")
	}
	user := &User{
		Id: userId,
	}
	_, err := o.Delete(user)
	return err
}
