package models

type Perm struct {
	Module string `orm:"size(20)"`
	Action string `orm:"size(20)"`
	Key    string `orm:"-"` // Module.Action
}
