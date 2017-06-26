package models

import (
	"errors"
	"time"
)

// 表结构
type Server struct {
	Id          int
	TypeId      int       // 0:普通服务器, 1:跳板机
	Ip          string    `orm:"size(20)"`  // 服务器IP
	Area        string    `orm:"size(20)"`  // 机房
	Description string    `orm:"size(200)"` // 服务器说明
	SshPort     int       // ssh端口
	SshUser     string    `orm:"size(50)"`                    // ssh用户
	SshPwd      string    `orm:"size(100)"`                   // ssh密码
	SshKey      string    `orm:"size(100)"`                   // ssh key路径
	WorkDir     string    `orm:"size(100)"`                   // 工作目录
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime  time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
}

const (
	SERVER_TYPE_NORMAL = 0 // 普通web服务器
	SERVER_TYPE_AGENT  = 1 // 跳板机
)

func (this *Server) table() string {
	return tableName("server")
}

func (this *Server) GetTotal(typeId int) (int64, error) {
	return o.QueryTable(this.table()).Filter("TypeId", typeId).Count()
}

// 获取一个服务器信息
func (this *Server) GetServer(id int, types ...int) (*Server, error) {
	var err error
	server := &Server{}
	server.Id = id
	if len(types) == 0 {
		err = o.Read(server)
	} else {
		err = o.QueryTable(this.table()).Filter("id", id).Filter("type_id", types[0]).One(server)
	}
	return server, err
}

// 获取跳板服务器列表
func (this *Server) GetAgentList(page, pageSize int) ([]Server, error) {
	var list []Server
	qs := o.QueryTable(this.table()).Filter("TypeId", SERVER_TYPE_AGENT)
	if pageSize > 0 {
		qs = qs.Limit(pageSize, (page-1)*pageSize)
	}
	_, err := qs.All(&list)
	return list, err
}

// 添加服务器
func (this *Server) AddServer(server *Server) error {
	server.Id = 0
	if o.Read(server, "ip"); server.Id > 0 {
		return errors.New("服务器IP已存在:" + server.Ip)
	}
	_, err := o.Insert(server)
	return err
}

// 根据id列表获取记录
func (this *Server) GetListByIds(ids []int) ([]Server, error) {
	var list []Server
	if len(ids) == 0 {
		return nil, errors.New("ids不能为空")
	}
	params := make([]interface{}, len(ids))
	for k, v := range ids {
		params[k] = v
	}
	_, err := o.QueryTable(this.table()).Filter("id__in", params...).All(&list)
	return list, err
}

// 获取普通服务器列表
func (this *Server) GetServerList(page, pageSize int) ([]Server, error) {
	var list []Server
	qs := o.QueryTable(this.table()).Filter("TypeId", SERVER_TYPE_NORMAL)
	if pageSize > 0 {
		qs = qs.Limit(pageSize, (page-1)*pageSize)
	}
	_, err := qs.All(&list)
	return list, err
}

// 修改服务器信息
func (this *Server) UpdateServer(server *Server, fields ...string) error {
	_, err := o.Update(server, fields...)
	return err
}

// 删除服务器
func (this *Server) DeleteServer(id int) error {
	_, err := o.QueryTable(this.table()).Filter("id", id).Delete()
	if err != nil {
		return err
	}
	return new(EnvServer).DeleteServer(id)
}
