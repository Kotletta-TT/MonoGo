package utils

import (
	"fmt"
	"reflect"
)

var floatType = reflect.TypeOf(float64(0))
var intType = reflect.TypeOf(int64(0))

func GetFloat64(val interface{}) float64 {
	v := reflect.ValueOf(val)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		panic(fmt.Sprintf("cannot convert %v to float64\n", v.Type()))
	}
	fv := v.Convert(floatType)
	return fv.Float()
}

func GetInt64(val interface{}) int64 {
	v := reflect.ValueOf(val)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(intType) {
		panic(fmt.Sprintf("cannot convert %v to int64\n", v.Type()))
	}
	fv := v.Convert(intType)
	return fv.Int()
}
