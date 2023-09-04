package nebulaModels

import (
	"errors"
	"fmt"
	"github.com/vesoft-inc/nebula-go/v3/nebula"
	"github.com/zhihu/norm/dialectors"
	"reflect"
)

var (
	NilPointError       = errors.New("assignment to entry in nil map")
	RecordNotFoundError = errors.New("record not found")
)

// UnmarshalResultSet 解组 ResultSet 为传入的结构体
func UnmarshalResultSet(resultSet *dialectors.ResultSet, in interface{}) error {
	switch values := in.(type) {
	case map[string]interface{}:
		return toMap(values, resultSet)
	case *map[string]interface{}:
		return toMap(*values, resultSet)
	case *[]map[string]interface{}:
		return toMapSlice(values, resultSet)
	default:
		val := reflect.ValueOf(values)
		switch val.Kind() {
		case reflect.Ptr:
			val = reflect.Indirect(val)
			switch val.Kind() {
			case reflect.Struct:
				return toStruct(val, resultSet)
			case reflect.Slice:
				return toStructSlice(val, resultSet)
			case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
				return toInt(val, resultSet)
			default:
				return errors.New(fmt.Sprintf("not support type. type is:%v", val.Kind()))
			}
		default:
			return errors.New("must be ptr")
		}
	}
}

func toInt(val reflect.Value, resultSet *dialectors.ResultSet) (err error) {
	if val.Interface() == nil {
		return NilPointError
	}

	if resultSet.GetRowSize() < 1 {
		val.SetInt(0)
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New("unknown exec error")
			}
		}
	}()
	cnt := resultSet.GetRows()[0].GetValues()[0].IVal
	val.SetInt(*cnt)
	return nil
}

func toStruct(val reflect.Value, resultSet *dialectors.ResultSet) (err error) {
	if val.Interface() == nil {
		return NilPointError
	}

	if resultSet.GetRowSize() < 1 {
		return RecordNotFoundError
	}

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New("unknown exec error")
			}
		}
	}()

	row := resultSet.GetRows()[0]
	fieldTagMap := getStructFieldTagMap(val.Type())
	for j, col := range resultSet.GetColNames() {
		fieldPos, ok := fieldTagMap[col]
		if !ok {
			continue
		}
		value := row.GetValues()[j]
		field := val.Field(fieldPos)
		err = setFieldValue(col, field, value)
	}

	return
}

func toStructSlice(val reflect.Value, resultSet *dialectors.ResultSet) (err error) {
	if val.Interface() == nil {
		return NilPointError
	}
	if resultSet.GetRowSize() < 1 {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New("unknown exec error")
			}
		}
	}()

	val.Set(reflect.MakeSlice(val.Type(), resultSet.GetRowSize(), resultSet.GetRowSize()))
	fieldTagMap := getStructFieldTagMap(val.Index(0).Type())
	for i, row := range resultSet.GetRows() {
		// 这里可以优化 GetColNames, 只循环两个共有的 key
		for j, col := range resultSet.GetColNames() {
			fieldPos, ok := fieldTagMap[col]
			if !ok {
				continue
			}
			nValue := row.GetValues()[j]
			field := val.Index(i).Field(fieldPos)
			err = setFieldValue(col, field, nValue)
		}
	}
	return
}

func toMap(values map[string]interface{}, resultSet *dialectors.ResultSet) error {
	if values == nil {
		return NilPointError
	}
	if resultSet.GetRowSize() < 1 {
		return RecordNotFoundError
	}
	row := resultSet.GetRows()[0]
	for i, col := range resultSet.GetColNames() {
		values[col] = nValueToInterface(row.Values[i])
	}
	return nil
}

func toMapSlice(values *[]map[string]interface{}, resultSet *dialectors.ResultSet) error {
	if values == nil {
		return NilPointError
	}

	cols := resultSet.GetColNames()
	_values := make([]map[string]interface{}, resultSet.GetRowSize())
	for i, row := range resultSet.GetRows() {
		_values[i] = make(map[string]interface{})
		for j, col := range cols {
			_values[i][col] = nValueToInterface(row.Values[j])
		}
	}
	*values = append(*values, _values...)
	return nil
}


func ParseNebulaResult(resList []map[string]interface{}) RespGraph {
	var res RespGraph
	resVertex := make([]*AddressTag, 0)
	resEdge := make([]*TxsEdge, 0)
	for _, res := range resList {
		var e TxsEdge
		for k, v := range res {

			if k == "total_amount" {
				e.TotalAmount = v.(float64)
			}
			if k == "tx_count" {
				e.TxCount = v.(int64)
			}

			//todo // map[string]interface{}
			//if mapRes, ok := v.(map[string]interface{}); ok {
			//	if _, ok := mapRes["v"]; ok {
			//		fmt.Println("v:",mapRes["v"])
			//	}
			//	if _, ok := mapRes["e"]; ok {
			//		fmt.Println("e:",mapRes["e"])
			//	}
			//
			//}
			// *nebula.Vertex
			if vertex, ok := v.(*nebula.Vertex); ok {
				tags := vertex.GetTags()
				if len(tags) == 0 {
					src := getValueofValue(vertex.GetVid())
					var a AddressTag
					a.Address = src
					resVertex = append(resVertex, &a)
				}else {
					srcs := getValueofAddresses(vertex.GetTags())
					for _, v := range srcs {
						resVertex = append(resVertex, v)
					}
				}
			}
			// *nebula.Edge
			if edge, ok := v.(*nebula.Edge); ok {
				e = getTxEdgeInfoFromProps(edge)
			}
			// *nebula.NList
			if nList, ok := v.(*nebula.NList); ok {
				pathValue := nList.GetValues()
				for _, v := range pathValue {
					if v.IsSetVVal() {
						for _, n := range v.GetVVal().GetTags() {
							src := getValueofAddress(n)
							resVertex = append(resVertex, &src)
						}
					}
					if v.IsSetEVal() {
						ed := getTxEdgeInfoFromProps(v.GetEVal())
						resEdge = append(resEdge, &ed)
					}
				}
			}
		}
		if !reflect.DeepEqual(e, TxsEdge{}) {
			resEdge = append(resEdge, &e)
		}
	}

	res.Edge = resEdge
	res.Vertex = resVertex
	return res
}


func ParseAddressTxNebula(resList []map[string]interface{}) RespGraph {
	var res RespGraph
	resVertex := make([]*AddressTag, 0)
	resEdge := make([]*TxsEdge, 0)

	for _, res := range resList {
		var e TxsEdge

		for k, v := range res {
			if k == "total_amount" {
				e.TotalAmount = v.(float64)
			}
			if k == "tx_count" {
				e.TxCount = v.(int64)
			}
		}
		resEdge = append(resEdge, &e)
	}
	res.Edge = resEdge
	res.Vertex = resVertex
	return res
}
