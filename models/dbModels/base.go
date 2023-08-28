package dbModels

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbhost := beego.AppConfig.String("mysql::dbhost")
	dbport := beego.AppConfig.String("mysql::dbport")
	dbuser := beego.AppConfig.String("mysql::dbuser")
	dbpassword := beego.AppConfig.String("mysql::dbpassword")
	dbname := beego.AppConfig.String("mysql::dbname")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	maxIdle := beego.AppConfig.DefaultInt("maxIdle", 16)  //连接池空闲
	maxConn := beego.AppConfig.DefaultInt("maxConn", 32) //连接池最大连接数  数据库默认链接数一般为512
	orm.RegisterDataBase("default", "mysql",
		dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=True",
		maxIdle, maxConn)
	// 初始化用户模块数据库
	userDBHost := beego.AppConfig.String("mysql::userdbhost")
	userDBPort := beego.AppConfig.String("mysql::userdbport")
	userDBUser := beego.AppConfig.String("mysql::userdbuser")
	userDBPassword := beego.AppConfig.String("mysql::userdbpassword")
	userDBName := beego.AppConfig.String("mysql::userdbname")
	orm.RegisterDataBase("user", "mysql",
		userDBUser+":"+userDBPassword+"@tcp("+userDBHost+":"+userDBPort+")/"+userDBName+"?charset=utf8&parseTime=True",
		5, 50)

	if "dev" == beego.AppConfig.String("runmode") {
		orm.Debug = true
	}
	beego.Info("orm init success")
}
