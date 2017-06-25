package models

import (
	"os"
	"time"
)

// 表结构
type Project struct {
	Id            int
	Name          string    `orm:"size(100)"`                   // 项目名称
	Domain        string    `orm:"size(100)"`                   // 项目标识
	Version       string    `orm:"size(20)"`                    // 最后发布版本
	VersionTime   time.Time `orm:"type(datetime)"`              // 最后发版时间
	RepoUrl       string    `orm:"size(100)"`                   // 仓库地址
	Status        int       `orm:"default(0)"`                  // 初始化状态
	ErrorMsg      string    `orm:"type(text)"`                  // 错误消息
	AgentId       int       `orm:"default(0)"`                  // 跳板机ID
	IgnoreList    string    `orm:"type(text)"`                  // 忽略文件列表
	BeforeShell   string    `orm:"type(text)"`                  // 发布前要执行的shell脚本
	AfterShell    string    `orm:"type(text)"`                  // 发布后要执行的shell脚本
	CreateVerfile int       `orm:"default(0)"`                  // 是否生成版本号文件
	VerfilePath   string    `orm:"size(50)"`                    // 版本号文件目录
	TaskReview    int       `orm:"default(0)"`                  // 发布单是否需要经过审批
	CreateTime    time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime    time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
}

func (this *Project) table() string {
	return tableName("project")
}

// 获取一个项目信息
func (this *Project) GetProject(id int) (*Project, error) {
	project := &Project{}
	project.Id = id
	if err := o.Read(project); err != nil {
		return nil, err
	}
	return project, nil
}

// 获取项目总数
func (this *Project) GetTotal() (int64, error) {
	return o.QueryTable(this.table()).Count()
}

// 获取项目列表
func (this *Project) GetList(page, pageSize int) ([]Project, error) {
	var list []Project
	offset := 0
	if pageSize == -1 {
		pageSize = 100000
	} else {
		offset = (page - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
	}

	_, err := o.QueryTable(this.table()).Offset(offset).Limit(pageSize).All(&list)
	return list, err
}

// 添加项目
func (this *Project) AddProject(project *Project) error {
	_, err := o.Insert(project)
	return err
}

// 更新项目信息
func (this *Project) UpdateProject(project *Project, fields ...string) error {
	_, err := o.Update(project, fields...)
	return err
}

// 克隆某个项目的仓库
func (this *Project) CloneRepo(projectId int) error {
	var reposityors Repository
	project, err := this.GetProject(projectId)
	if err != nil {
		return err
	}

	err = reposityors.CloneRepo(project.RepoUrl, GetProjectPath(project.Domain))
	if err != nil {
		project.Status = -1
		project.ErrorMsg = err.Error()
	} else {
		project.Status = 1
	}
	this.UpdateProject(project, "Status", "ErrorMsg")

	return err
}

// 删除一个项目
func (this *Project) DeleteProject(projectId int) error {
	var (
		tasks Task
		envs  Env
	)
	project, err := this.GetProject(projectId)
	if err != nil {
		return err
	}
	// 删除目录
	path := GetProjectPath(project.Domain)
	os.RemoveAll(path)
	// 环境配置
	if envList, err := envs.GetEnvListByProjectId(project.Id); err != nil {
		for _, env := range envList {
			envs.DeleteEnv(env.Id)
		}
	}
	// 删除任务
	tasks.DeleteByProjectId(project.Id)
	// 删除项目
	o.Delete(project)
	return nil
}
