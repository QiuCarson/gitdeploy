package models

import (
	"errors"
	"time"
)

// 角色
type Role struct {
	Id          int
	RoleName    string     `orm:"size(20)"`                    // 角色名称
	ProjectIds  string     `orm:"size(1000)"`                  // 项目权限
	Description string     `orm:"size(200)"`                   // 说明
	CreateTime  time.Time  `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime  time.Time  `orm:"auto_now;type(datetime)"`     // 更新时间
	PermList    []Roleperm `orm:"-"`                           // 权限列表
	UserList    []User     `orm:"-"`                           // 用户列表
}

func (this *Role) table() string {
	return tableName("role")
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

type Roleperm struct {
	Id      int
	Module  string `orm:"size(20)"`
	Action  string `orm:"size(20)"`
	Keyinfo string `orm:"-"` // Module.Action
}

func (this *Role) loadRoleExtra(role *Role) {
	sql := "SELECT SUBSTRING_INDEX(perm, '.', 1) as module,SUBSTRING_INDEX(perm, '.', -1) as `action` , perm AS `keyinfo`FROM " + tableName("role_perm") + " WHERE role_id = ?"
	o.Raw(sql, role.Id).QueryRows(&role.PermList)
}

// 添加角色
func (this *Role) AddRole(role *Role) error {
	if _, err := this.GetRoleByName(role.RoleName); err == nil {
		return errors.New("角色已存在")
	}
	_, err := o.Insert(role)
	return err
}

// 根据名称获取角色
func (this *Role) GetRoleByName(roleName string) (*Role, error) {
	role := &Role{
		RoleName: roleName,
	}
	if err := o.Read(role, "RoleName"); err != nil {
		return nil, err
	}
	this.loadRoleExtra(role)
	return role, nil
}

// 更新角色信息
func (this *Role) UpdateRole(role *Role, fields ...string) error {
	if v, err := this.GetRoleByName(role.RoleName); err == nil && v.Id != role.Id {
		return errors.New("角色名称已存在")
	}
	_, err := o.Update(role, fields...)
	return err
}

// 设置角色权限
func (this *Role) SetPerm(roleId int, perms []string) error {
	if _, err := this.GetRole(roleId); err != nil {
		return err
	}
	all := new(SystemService).GetPermList()
	pmmap := make(map[string]bool)
	for _, list := range all {
		for _, perm := range list {
			pmmap[perm.Keyinfo] = true
		}
	}
	for _, v := range perms {
		if _, ok := pmmap[v]; !ok {
			return errors.New("权限名称无效:" + v)
		}
	}
	o.Raw("DELETE FROM "+tableName("role_perm")+" WHERE role_id = ?", roleId).Exec()
	for _, v := range perms {
		o.Raw("REPLACE INTO "+tableName("role_perm")+" (role_id, perm) VALUES (?, ?)", roleId, v).Exec()
	}
	return nil
}

// 删除角色
func (this *Role) DeleteRole(id int) error {
	role, err := this.GetRole(id)
	if err != nil {
		return err
	}
	o.Delete(role)
	o.Raw("DELETE FROM "+tableName("role_user")+" WHERE role_id = ?", id).Exec()
	return nil
}

// 获取所有角色列表
func (this *Role) GetAllRoles() ([]Role, error) {
	var (
		roles []Role // 角色列表
	)
	if _, err := o.QueryTable(this.table()).All(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}
