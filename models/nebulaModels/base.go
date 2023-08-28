package nebulaModels

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/vesoft-inc/nebula-go/v3/nebula"
	"github.com/zhihu/norm"
	"github.com/zhihu/norm/dialectors"
	"sync"
	"time"
)

type Config struct {
	DbHost     string `json:"db_host"`
	DbSpace    string `json:"db_space"`
	DbUser     string `json:"db_username"`
	DbPassword string `json:"db_password"`
}

var nebulaConfig *Config

func init() {
	nebulaConfig = &Config{
		DbHost:     beego.AppConfig.String("nebula::dbhost"),
		DbSpace:    beego.AppConfig.String("nebula::dbspace"),
		DbUser:     beego.AppConfig.String("nebula::dbuser"),
		DbPassword: beego.AppConfig.String("nebula::dbpassword"),
	}
	beego.Info("nebula config init success")
}

func newDB() *norm.DB {
	dialector := dialectors.MustNewNebulaDialector(dialectors.DialectorConfig{
		Addresses: []string{nebulaConfig.DbHost},
		Timeout:   time.Second * 5,
		Space:     nebulaConfig.DbSpace,
		Username:  nebulaConfig.DbUser,
		Password:  nebulaConfig.DbPassword,
	})
	db := norm.MustOpen(dialector, norm.Config{})
	return db
}

func Init() *norm.DB {
	var nebulaDB *norm.DB
	var once sync.Once
	once.Do(func() {
		nebulaDB = newDB()
	})
	return nebulaDB
}

func getValueofValue(value *nebula.Value) string {
	if value.NVal != nil {
		return fmt.Sprintf("value.NVal=%v", value.NVal)
	}
	if value.BVal != nil {
		return fmt.Sprintf("value.BVal=%v", value.BVal)
	}
	if value.IVal != nil {
		return fmt.Sprintf("value.IVal=%v", value.IVal)
	}
	if value.FVal != nil {
		return fmt.Sprintf("value.FVal=%v", value.FVal)
	}
	if value.SVal != nil {
		return fmt.Sprintf("%v", string(value.SVal))
	}
	if value.DVal != nil {
		return fmt.Sprintf("value.DVal=%v", value.DVal)
	}
	if value.TVal != nil {
		return fmt.Sprintf("value.TVal=%v", value.TVal)
	}
	if value.DtVal != nil {
		return fmt.Sprintf("value.DtVal=%v", value.DtVal)
	}
	if value.VVal != nil {
		return fmt.Sprintf("value.VVal=%v", value.VVal)
	}
	if value.EVal != nil {
		return fmt.Sprintf("value.EVal=%v", value.EVal)
	}
	if value.PVal != nil {
		return fmt.Sprintf("value.PVal=%v", value.PVal)
	}
	if value.LVal != nil {
		return fmt.Sprintf("value.LVal=%v", value.LVal)
	}
	if value.MVal != nil {
		return fmt.Sprintf("value.MVal=%v", value.MVal)
	}
	if value.UVal != nil {
		return fmt.Sprintf("value.UVal=%v", value.UVal)
	}
	if value.GVal != nil {
		return fmt.Sprintf("value.GVal=%v", value.GVal)
	}
	if value.GgVal != nil {
		return fmt.Sprintf("value.GgVal=%v", value.GgVal)
	}
	if value.DuVal != nil {
		return fmt.Sprintf("value.DuVal=%v", value.DuVal)
	}
	return ""
}
