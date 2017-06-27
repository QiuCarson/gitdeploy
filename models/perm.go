package models

type Perm struct {
	Id      int
	Module  string `orm:"size(20)"`
	Action  string `orm:"size(20)"`
	Keyinfo string `orm:"-"` // Module.Action
}
