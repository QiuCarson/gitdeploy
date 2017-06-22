package models

import "time"

// 角色
type Role struct {
	Id          int
	RoleName    string    `orm:"size(20)"`                    // 角色名称
	ProjectIds  string    `orm:"size(1000)"`                  // 项目权限
	Description string    `orm:"size(200)"`                   // 说明
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime  time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	PermList    []Perm    `orm:"-"`                           // 权限列表
	UserList    []User    `orm:"-"`                           // 用户列表
}

// 根据id获取角色信息
func (this *Role) GetRole(id int) (*Role, error) {
	role := &Role{
		Id: id,
	}
	err := o.Read(role)
	if err != nil {
		return nil, err
	}
	this.loadRoleExtra(role)
	return role, err
}
func (this *Role) loadRoleExtra(role *Role) {
	o.Raw("SELECT SUBSTRING_INDEX(perm, '.', 1) as module,SUBSTRING_INDEX(perm, '.', -1) as `action`, perm AS `key` FROM "+tableName("role_perm")+" WHERE role_id = ?", role.Id).QueryRows(&role.PermList)
}
