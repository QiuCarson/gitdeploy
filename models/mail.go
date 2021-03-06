package models

import (
	"time"
)

// 表结构
type MailTpl struct {
	Id         int
	UserId     int
	Name       string    `orm:"size(100)"`                   // 模板名
	Subject    string    `orm:"size(200)"`                   // 邮件主题
	Content    string    `orm:"type(text)"`                  // 邮件内容
	MailTo     string    `orm:"size(1000)"`                  // 预设收件人
	MailCc     string    `orm:"size(1000)"`                  // 预设抄送人
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
}

func (this *MailTpl) table() string {
	return tableName("mail_tpl")
}

// 获取邮件模板列表
func (this *MailTpl) GetMailTplList() ([]MailTpl, error) {
	var list []MailTpl
	_, err := o.QueryTable(this.table()).OrderBy("-id").All(&list)
	return list, err
}

func (this *MailTpl) AddMailTpl(tpl *MailTpl) error {
	_, err := o.Insert(tpl)
	return err
}

func (this *MailTpl) GetMailTpl(id int) (*MailTpl, error) {
	tpl := &MailTpl{}
	tpl.Id = id
	err := o.Read(tpl)
	return tpl, err
}

func (this *MailTpl) SaveMailTpl(tpl *MailTpl) error {
	_, err := o.Update(tpl)
	return err
}

func (this *MailTpl) DelMailTpl(id int) error {
	_, err := o.QueryTable(this.table()).Filter("id", id).Delete()
	return err
}
