package models

import (
	"time"

	"github.com/astaxie/beego/orm"
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
func (this *EnvServer) serverTable() string {
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
	_, err := o.QueryTable(new(EnvServer).serverTable()).Filter("env_id", envId).All(&list)
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
	o.QueryTable(new(EnvServer).serverTable()).Filter("env_id", id).Delete()
	return nil
}

// 新增发布环境
func (this *EnvServer) AddEnv(env *Env) error {
	env.ServerCount = len(env.ServerList)
	if _, err := o.Insert(env); err != nil {
		return err
	}
	for _, sv := range env.ServerList {
		es := new(EnvServer)
		es.ProjectId = env.ProjectId
		es.EnvId = env.Id
		es.ServerId = sv.Id
		o.Insert(es)
	}
	return nil
}

// 保存环境配置
func (this *EnvServer) SaveEnv(env *Env) error {
	env.ServerCount = len(env.ServerList)
	if _, err := o.Update(env); err != nil {
		return err
	}
	o.QueryTable(this.serverTable()).Filter("env_id", env.Id).Delete()
	for _, sv := range env.ServerList {
		es := new(EnvServer)
		es.ProjectId = env.ProjectId
		es.EnvId = env.Id
		es.ServerId = sv.Id
		o.Insert(es)
	}
	return nil
}

// 删除发布环境
func (this *EnvServer) DeleteEnv(id int) error {
	o.QueryTable(new(Env).table()).Filter("id", id).Delete()
	o.QueryTable(this.serverTable()).Filter("env_id", id).Delete()
	return nil
}

// 获取一个发布环境信息
func (this *EnvServer) GetEnv(id int) (*Env, error) {
	env := &Env{}
	env.Id = id
	err := o.Read(env)
	if err == nil {
		env.ServerList, _ = new(Env).GetEnvServers(env.Id)
	}
	return env, err
}

// 删除服务器
func (this *EnvServer) DeleteServer(serverId int) error {
	var envServers []EnvServer
	o.QueryTable(this.serverTable()).Filter("server_id", serverId).All(&envServers)
	if len(envServers) < 1 {
		return nil
	}
	envIds := make([]int, 0, len(envServers))
	for _, v := range envServers {
		envIds = append(envIds, v.EnvId)
	}
	o.QueryTable(this.serverTable()).Filter("server_id", serverId).Delete()
	o.QueryTable(new(Env).table()).Filter("id__in", envIds).Update(orm.Params{
		"server_count": orm.ColValue(orm.ColMinus, 1),
	})
	return nil
}

// 根据服务器id发布环境列表
func (this *EnvServer) GetEnvListByServerId(serverId int) ([]Env, error) {
	var (
		servList []EnvServer
		envList  []Env
	)
	o.QueryTable(this.serverTable()).Filter("server_id", serverId).All(&servList)
	envIds := make([]int, 0, len(servList))
	for _, serv := range servList {
		envIds = append(envIds, serv.EnvId)
	}
	envList = make([]Env, 0)
	if len(envIds) > 0 {
		if _, err := o.QueryTable(new(Env).table()).Filter("id__in", envIds).All(&envList); err != nil {
			return envList, err
		}
	}
	return envList, nil
}
