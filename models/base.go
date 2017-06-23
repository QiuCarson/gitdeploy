package models

import (
	"crypto/md5"
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	o           orm.Ormer
	tablePrefix string
)

func init() {
	dbHost := beego.AppConfig.String("db.host")
	dbPort := beego.AppConfig.String("db.port")
	dbUser := beego.AppConfig.String("db.user")
	dbPassword := beego.AppConfig.String("db.password")
	dbName := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	tablePrefix = beego.AppConfig.String("db.prefix")

	if dbPort == "" {
		dbPort = "3306"
	}
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"
	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModelWithPrefix(tablePrefix,
		new(User),
		new(Role),
		new(Perm),
		new(Action),
		new(Server),
		new(Project),
		new(Task),
	)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	o = orm.NewOrm()
	//orm.RunCommand()

}

// 返回真实表名
func tableName(name string) string {
	return tablePrefix + name
}

// 生成md5
func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func DBVersion() string {
	var lists []orm.ParamsList
	o.Raw("SELECT VERSION()").ValuesList(&lists)
	return lists[0][0].(string)
}
