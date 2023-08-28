package nebulaModels

import (
	"reflect"
	"time"

	nebula_type "github.com/vesoft-inc/nebula-go/v3/nebula"
	"github.com/zhihu/norm/constants"
)

//var rankMap map[string]int64
//
//func isRankMapEmpty() bool {
//	if len(rankMap) == 0 {
//		return true
//	}
//	return false
//}
//
//func GetRank(key string) (int64, bool) {
//	rankMap = make(map[string]int64)
//	value, ok := rankMap[key]
//	if ok {
//		rankMap[key] = value + 1
//		return value, ok
//	} else {
//		return 0, ok
//	}
//}
//
//func SetRank(key string, txRank int64) {
//	rankMap = make(map[string]int64)
//	rankMap[key] = txRank
//}

// getStructFieldTagMap 将 struct 中标记为 norm 的 tag 提取出来, 并记录 field 的位置
func getStructFieldTagMap(typ reflect.Type) map[string]int {
	tagMap := make(map[string]int)
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get(constants.StructTagName)
		if tag == "" || tag == "-" {
			continue
		}
		tagMap[tag] = i
	}
	return tagMap
}

// setFieldValue 将 nvalue 的值设置到 struct.field 上, 并自动转换类型
func setFieldValue(tag string, field reflect.Value, nValue *nebula_type.Value) error {
	switch field.Kind() {
	case reflect.Bool:
		field.SetBool(nValue.GetBVal())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(nValue.GetIVal())
	case reflect.Float32, reflect.Float64:
		field.SetFloat(nValue.GetFVal())
	case reflect.String:
		field.SetString(string(nValue.GetSVal()))
	case reflect.Struct:
		switch field.Type().String() {
		case "time.Time":
			ts := nValue.GetIVal()
			field.Set(reflect.ValueOf(time.Unix(ts, 0)))
		default:
			//fmt.Printf("debug: type[%v] mapping not implement\n", field.Type().String())
		}
	default:
		//fmt.Printf("debug: type[%v] mapping not implement\n", field.Type().String())
		return nil
	}
	return nil
}

// nValueToInterface 将 nvalue 的值转换类型并返回interface
func nValueToInterface(p *nebula_type.Value) interface{} {
	if p.IsSetNVal() {
		return p.GetNVal()
	}
	if p.IsSetBVal() {
		return p.GetBVal()
	}
	if p.IsSetIVal() {
		return p.GetIVal()
	}
	if p.IsSetFVal() {
		return p.GetFVal()
	}
	if p.IsSetSVal() {
		return p.GetSVal()
	}
	if p.IsSetDVal() {
		return p.GetDVal()
	}
	if p.IsSetTVal() {
		return p.GetTVal()
	}
	if p.IsSetDtVal() {
		return p.GetDtVal()
	}
	if p.IsSetVVal() {
		return p.GetVVal()
	}
	if p.IsSetEVal() {
		return p.GetEVal()
	}
	if p.IsSetPVal() {
		return p.GetPVal()
	}
	if p.IsSetLVal() {
		return p.GetLVal()
	}
	if p.IsSetMVal() {
		return p.GetMVal()
	}
	if p.IsSetUVal() {
		return p.GetUVal()
	}
	if p.IsSetGVal() {
		return p.GetGVal()
	}
	if p.IsSetGgVal() {
		return p.GetGgVal()
	}
	if p.IsSetDuVal() {
		return p.GetDuVal()
	}
	return nil
}
