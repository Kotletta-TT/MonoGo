package utils

import (
	"reflect"
	"runtime"
)

var floatType = reflect.TypeOf(float64(0))
var intType = reflect.TypeOf(int64(0))

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
