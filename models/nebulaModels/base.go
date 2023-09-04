package nebulaModels

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/server/txs-analysis/constant"
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

// ------------------解析边-----------------

func getValueofValue(value *nebula.Value) string {
	if value.NVal != nil {
		return fmt.Sprintf("%v", value.NVal)
	}
	if value.BVal != nil {
		return fmt.Sprintf("%v", value.BVal)
	}
	if value.SVal != nil {
		return fmt.Sprintf("%v", string(value.SVal))
	}
	if value.DVal != nil {
		return fmt.Sprintf("%v", value.DVal)
	}
	if value.TVal != nil {
		return fmt.Sprintf("%v", value.TVal)
	}
	if value.DtVal != nil {
		return value.DtVal.String()
	}
	if value.VVal != nil {
		return fmt.Sprintf("%v", value.VVal)
	}
	if value.EVal != nil {
		return fmt.Sprintf("%v", value.EVal)
	}
	if value.PVal != nil {
		return fmt.Sprintf("%v", value.PVal)
	}
	if value.LVal != nil {
		return fmt.Sprintf("%v", value.LVal)
	}
	if value.MVal != nil {
		return fmt.Sprintf("%v", value.MVal)
	}
	if value.UVal != nil {
		return fmt.Sprintf("%v", value.UVal)
	}
	if value.GVal != nil {
		return fmt.Sprintf("%v", value.GVal)
	}
	if value.GgVal != nil {
		return fmt.Sprintf("%v", value.GgVal)
	}
	if value.DuVal != nil {
		return fmt.Sprintf("%v", value.DuVal)
	}
	return ""
}

func getFloatValueofValue(value *nebula.Value) float64 {
	if value.FVal != nil {
		return *value.FVal
	}
	return 0
}

func getDtValueofValue(value *nebula.Value) nebula.DateTime {
	if value.DtVal != nil {
		return *value.DtVal
	}
	return nebula.DateTime{}
}

func getIntValueofValue(value *nebula.Value) int64 {
	if value.IVal != nil {
		return *value.IVal
	}
	return 0
}

// ------------------解析点-----------------

//func getValueofTag(tag *nebula.Tag) string {
//	v := ""
//	s := fmt.Sprintf("name:%s", string(tag.Name))
//	v += s
//	for k, p := range tag.Props {
//		s = fmt.Sprintf("[%s]=%s\n", k, getValueofValue(p))
//		v += s
//	}
//	return v
//}

func getValueofAddress(tag *nebula.Tag) AddressTag {
	var address AddressTag
	for _, v := range tag.Props {
		if v.IVal != nil {
			address.Type = int(*v.IVal)
		}
		if v.SVal != nil {
			address.Address = string(v.SVal)
		}
	}
	return address
}

func getValueofAddresses(tags []*nebula.Tag) []*AddressTag {
	var addresses []*AddressTag
	for _, tag := range tags {
		tag := getValueofAddress(tag)
		addresses = append(addresses, &tag)
	}
	return addresses
}

//func getValueofTags(tags []*nebula.Tag) string {
//	v := ""
//	for i, tag := range tags {
//		s := fmt.Sprintf("t[%d]=%s\n", i, getValueofTag(tag))
//		v += s
//	}
//	return v
//}

// getTxEdgeInfoFromProps 获取tx中的值
func getTxEdgeInfoFromProps(edge *nebula.Edge) TxsEdge {
	tx := TxsEdge{}
	if value, exist := edge.Props[constant.BlockId]; exist {
		tx.BlockId = getIntValueofValue(value)
	}
	if value, exist := edge.Props[constant.TxHash]; exist {
		tx.TxHash = getValueofValue(value)
	}
	if txTime, exist := edge.Props[constant.TxTime]; exist {
		tx.TxTime = getValueofValue(txTime)
	}
	if fromAddress, exist := edge.Props[constant.FromAddress]; exist {
		tx.FromAddress = getValueofValue(fromAddress)
	}
	if toAddress, exist := edge.Props[constant.ToAddress]; exist {
		tx.ToAddress = getValueofValue(toAddress)
	}
	if value, exist := edge.Props[constant.Value]; exist {
		tx.Value = getValueofValue(value)
	}
	if amount, exist := edge.Props[constant.Amount]; exist {
		tx.Amount = getFloatValueofValue(amount)
	}
	if value, exist := edge.Props[constant.TokenId]; exist {
		tx.TokenId = getValueofValue(value)
	}
	if value, exist := edge.Props[constant.TxType]; exist {
		tx.TxType = getIntValueofValue(value)
	}
	if value, exist := edge.Props[constant.TokenAddress]; exist {
		tx.TokenAddress = getValueofValue(value)
	}
	return tx
}
