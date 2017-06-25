package models

import (
	"time"
)

// 发布环境
type Env struct {
	Id          int
	ProjectId   int       `orm:"index"`                       // 项目id
	Name        string    `orm:"size(20)"`                    // 发布环境名称
	SshUser     string    `orm:"size(20)"`                    // 发布帐号
	SshPort     string    `orm:"size(10)"`                    // SSH端口
	SshKey      string    `orm:"size(100)"`                   // SSH KEY路径
	PubDir      string    `orm:"size(100)"`                   // 发布目录
	BeforeShell string    `orm:"type(text)"`                  // 发布前执行的shell脚本
	AfterShell  string    `orm:"type(text)"`                  // 发布后执行的shell脚本
	ServerCount int       `orm:"default(0)"`                  // 服务器数量
	SendMail    int       `orm:"default(0)"`                  // 是否发送发版邮件通知
	MailTplId   int       `orm:"default(0)"`                  // 邮件模板id
	MailTo      string    `orm:"size(1000)"`                  // 邮件收件人
	MailCc      string    `orm:"size(1000)"`                  // 邮件抄送人
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime  time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	ServerList  []Server  `orm:"-"`                           // 服务器列表
}

// 表结构
type EnvServer struct {
	Id        int
	ProjectId int `orm:"default(0)"`       // 项目id
	EnvId     int `orm:"default(0);index"` // 环境id
	ServerId  int `orm:"default(0)"`       // 服务器id
}

func (this *Env) table() string {
	return tableName("env")
}
func (this *Env) serverTable() string {
	return tableName("env_server")
}

// 获取某个项目的发布环境列表
func (this *Env) GetEnvListByProjectId(projectId int) ([]Env, error) {
	var list []Env
	_, err := o.QueryTable(this.table()).Filter("project_id", projectId).All(&list)
	for _, env := range list {
		env.ServerList, _ = this.GetEnvServers(env.Id)
	}
	return list, err
}

// 获取某个发布环境的服务器列表
func (this *Env) GetEnvServers(envId int) ([]Server, error) {
	var (
		list    []EnvServer
		servers Server
	)
	_, err := o.QueryTable(this.serverTable()).Filter("env_id", envId).All(&list)
	if err != nil {
		return nil, err
	}
	servIds := make([]int, 0, len(list))
	for _, v := range list {
		servIds = append(servIds, v.ServerId)
	}

	return servers.GetListByIds(servIds)
}

// 删除发布环境
func (this *Env) DeleteEnv(id int) error {
	o.QueryTable(this.table()).Filter("id", id).Delete()
	o.QueryTable(this.serverTable()).Filter("env_id", id).Delete()
	return nil
}
