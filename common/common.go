package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// 将数据库中查询到一条数据库中的每个字段，对应填写到model结构体中
// data : 数据库中查询出来的一条数据，查询结果用map表示，Key是列名 value是对应的值
// model : Server用来接收data的Struct
// func DataToStructByTagSql(data map[string]string, model interface{}) {
// 	valModel := reflect.ValueOf(model).Elem()
// 	typModel := reflect.TypeOf(model).Elem()
// 	// 依次遍历ModelStruct的每个字段，将对应的data的字段填入
// 	num := valModel.NumField()
// 	for i := 0; i < num; i++ {
// 		// 1. 通过Model字段对应的Tag SQL找到data对应的值
// 		dataValue := data[typModel.Field(i).Tag.Get("sql")]
// 		// 2. Model字段对应名字
// 		nameModel := typModel.Field(i).Name //比如Name Age， Field()直接显示这些字段对应的值
// 		// 3. Model字段类型 & SQL查询结果类型是否匹配
// 		typeModel := valModel.Field(i).Type() //比如string、int对应的reflect.Type
// 		dataValRef := reflect.ValueOf(dataValue)
// 		if typeModel != dataValRef.Type() {
// 			// 类型转换：将data转换成Model字段对应的类型的reflect的类型
// 			TypeConversion(dataValue, typeModel.Name()) //。Name可以将Type转换成传统类型
// 		}
// 		// 4. 填入Model中
// 		valModel.FieldByName(nameModel).Set(dataValRef)
// 	}
// }
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		//获取sql对应的值
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		//获取对应字段的名称
		name := objValue.Type().Field(i).Name
		//获取对应字段类型
		structFieldType := objValue.Field(i).Type()
		//获取变量类型，也可以直接写"string类型"
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			//类型转换
			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
			if err != nil {

			}
		}
		//设置类型值
		objValue.FieldByName(name).Set(val)
	}
}

//类型转换：将value转换成ntype类型，并返回value转换后对应reflect.Value类型的结果值
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
