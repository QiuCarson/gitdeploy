package models

import (
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Task struct {
	Id           int
	ProjectId    int       `orm:"index"`                       // 项目id
	StartVer     string    `orm:"size(20)"`                    // 起始版本号
	EndVer       string    `orm:"size(20)"`                    // 结束版本号
	Message      string    `orm:"type(text)"`                  // 版本说明
	UserId       int       `orm:"index"`                       // 创建人ID
	UserName     string    `orm:"size(20)"`                    // 创建人名称
	BuildStatus  int       `orm:"default(0)"`                  // 构建状态
	ChangeLogs   string    `orm:"type(text)"`                  // 修改日志列表
	ChangeFiles  string    `orm:"type(text)"`                  // 修改文件列表
	Filepath     string    `orm:"size(200)"`                   // 更新包路径
	PubEnvId     int       `orm:"default(0)"`                  // 发布环境ID
	PubStatus    int       `orm:"default(0)"`                  // 发布状态：1 正在发布，2 发布到跳板机，3 发布到目标服务器，-2 发布到跳板机失败，-3 发布到目标服务器失败
	PubTime      time.Time `orm:"null;type(datetime)"`         // 发布时间
	ErrorMsg     string    `orm:"type(text)"`                  // 错误消息
	PubLog       string    `orm:"type(text)"`                  // 发布日志
	ReviewStatus int       `orm:"default(0)"`                  // 审批状态
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime   time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	ProjectInfo  *Project  `orm:"-"`                           // 项目信息
	EnvInfo      *Env      `orm:"-"`                           // 发布环境
}

func (this *Task) table() string {
	return tableName("task")
}

func (this *Task) GetProjectPubStat() []map[string]int {
	var maps []orm.Params
	sql := "SELECT project_id, COUNT(*) AS count FROM " + this.table() + " WHERE pub_status = 3 GROUP BY project_id ORDER BY `count` DESC"
	num, err := o.Raw(sql).Values(&maps)
	result := make([]map[string]int, 0, num)
	if err == nil && num > 0 {
		for _, v := range maps {
			projectId, _ := strconv.Atoi(v["project_id"].(string))
			count, _ := strconv.Atoi(v["count"].(string))
			result = append(result, map[string]int{
				"project_id": projectId,
				"count":      count,
			})
		}
	}
	return result
}

// 获取一个任务信息
func (this *Task) GetTask(id int) (*Task, error) {
	var p Project
	task := &Task{}
	task.Id = id
	if err := o.Read(task); err != nil {
		return nil, err
	}
	task.ProjectInfo, _ = p.GetProject(task.ProjectId)
	return task, nil
}

// 获取已发布任务总数
func (this *Task) GetPubTotal() (int64, error) {
	return o.QueryTable(this.table()).Filter("pub_status", 3).Count()
}
