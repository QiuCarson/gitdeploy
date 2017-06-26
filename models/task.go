package models

import (
	"fmt"
	"os"
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

type TaskReview struct {
	Id         int
	TaskId     int       `orm:"default(0)"`                  // 任务id
	UserId     int       `orm:"default(0)"`                  // 审批人id
	UserName   string    `orm:"size(20)"`                    // 审批人
	Status     int       `orm:"default(0)"`                  // 审批结果(1:通过;0:不通过)
	Message    string    `orm:"type(text)"`                  // 审批说明
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
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

// 删除某个项目下的而所有发布任务
func (this *Task) DeleteByProjectId(projectId int) error {
	_, err := o.QueryTable(this.table()).Filter("project_id", projectId).Delete()
	return err
}

// 删除任务
func (this *Task) DeleteTask(taskId int) error {
	task, err := this.GetTask(taskId)
	if err != nil {
		return err
	}
	if _, err := o.Delete(task); err != nil {
		return err
	}
	return os.RemoveAll(GetTaskPath(task.Id))
}

// 获取任务单列表
func (this *Task) GetList(page, pageSize int, filters ...interface{}) ([]Task, int64) {
	var (
		list  []Task
		count int64
	)

	offset := (page - 1) * pageSize
	query := o.QueryTable(this.table())

	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			field, ok := filters[k].(string)
			if !ok {
				continue
			}
			switch field {
			case "start_date":
				v := fmt.Sprintf("%s 00:00:00", filters[k+1].(string))
				query = query.Filter("create_time__gte", v)
			case "end_date":
				v := fmt.Sprintf("%s 23:59:59", filters[k+1].(string))
				query = query.Filter("create_time__lte", v)
			default:
				v := filters[k+1]
				query = query.Filter(filters[k].(string), v)
			}
		}
	}
	count, _ = query.Count()

	if count > 0 {
		query.OrderBy("-id").Offset(offset).Limit(pageSize).All(&list)
		for k, v := range list {
			if p, err := new(Project).GetProject(v.ProjectId); err == nil {
				list[k].ProjectInfo = p
			} else {
				list[k].ProjectInfo = new(Project)
			}
		}
	}

	return list, count
}

// 添加任务
func (this *Task) AddTask(task *Task) error {
	if _, err := new(EnvServer).GetEnv(task.PubEnvId); err != nil {
		return fmt.Errorf("获取环境信息失败: %s", err.Error())
	}
	project, err := new(Project).GetProject(task.ProjectId)
	if err != nil {
		return fmt.Errorf("获取项目信息失败: %s", err.Error())
	}
	if project.TaskReview > 0 {
		task.ReviewStatus = 0 // 未审批
	} else {
		task.ReviewStatus = 1 // 已审批
	}
	task.PubStatus = 0
	// task.PubTime = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	_, err = o.Insert(task)
	return err
}

// 构建发布包
func (this *Task) BuildTask(task *Task) {
	err := new(Deploy).Build(task)
	if err != nil {
		task.BuildStatus = -1
		task.ErrorMsg = err.Error()
	} else {
		task.BuildStatus = 1
		task.ErrorMsg = ""
	}
	this.UpdateTask(task, "BuildStatus", "ErrorMsg")
}

// 更新任务信息
func (this *Task) UpdateTask(task *Task, fields ...string) error {
	_, err := o.Update(task, fields...)
	return err
}

// 获取审批信息
func (this *Task) GetReviewInfo(taskId int) (*TaskReview, error) {
	review := new(TaskReview)
	err := o.QueryTable(tableName("task_review")).Filter("task_id", taskId).OrderBy("-id").Limit(1).One(review)
	return review, err
}

// 任务审批
func (this *Task) ReviewTask(taskId, userId, status int, message string) error {
	if status != 1 && status != -1 {
		return fmt.Errorf("审批状态无效: %d", status)
	}
	user, err := new(User).GetUser(userId, false)
	if err != nil {
		return err
	}
	task, err := this.GetTask(taskId)
	if err != nil {
		return err
	}
	review := &TaskReview{}
	review.TaskId = task.Id
	review.UserId = user.Id
	review.UserName = user.UserName
	review.Status = status
	review.Message = message
	if _, err := o.Insert(review); err != nil {
		return err
	}

	task.ReviewStatus = status
	return this.UpdateTask(task, "ReviewStatus")
}

// 发布统计
func (this *Task) GetPubStat(rangeType string) map[int]int {
	var sql string
	var maps []orm.Params

	switch rangeType {
	case "this_month":
		year, month, _ := time.Now().Date()
		startTime := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTime := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		sql = fmt.Sprintf("SELECT DAY(pub_time) AS date, COUNT(*) AS count FROM %s WHERE pub_time BETWEEN '%s' AND '%s' GROUP BY DAY(pub_time) ORDER BY `date` ASC", this.table(), startTime, endTime)
	case "last_month":
		year, month, _ := time.Now().AddDate(0, -1, 0).Date()
		startTime := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTime := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		sql = fmt.Sprintf("SELECT DAY(pub_time) AS date, COUNT(*) AS count FROM %s WHERE pub_time BETWEEN '%s' AND '%s' GROUP BY DAY(pub_time) ORDER BY `date` ASC", this.table(), startTime, endTime)
	case "this_year":
		year := time.Now().Year()
		startTime := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTime := fmt.Sprintf("%d-12-31 23:59:59", year)
		sql = fmt.Sprintf("SELECT MONTH(pub_time) AS date, COUNT(*) AS count FROM %s WHERE pub_time BETWEEN '%s' AND '%s' GROUP BY MONTH(pub_time) ORDER BY `date` ASC", this.table(), startTime, endTime)
	case "last_year":
		year := time.Now().Year() - 1
		startTime := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTime := fmt.Sprintf("%d-12-31 23:59:59", year)
		sql = fmt.Sprintf("SELECT MONTH(pub_time) AS date, COUNT(*) AS count FROM %s WHERE pub_time BETWEEN '%s' AND '%s' GROUP BY MONTH(pub_time) ORDER BY `date` ASC", this.table(), startTime, endTime)
	}

	num, err := o.Raw(sql).Values(&maps)

	result := make(map[int]int)
	if err == nil && num > 0 {
		for _, v := range maps {
			date, _ := strconv.Atoi(v["date"].(string))
			count, _ := strconv.Atoi(v["count"].(string))
			result[date] = count
		}
	}
	return result
}
