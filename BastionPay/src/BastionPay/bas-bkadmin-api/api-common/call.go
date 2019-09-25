package common

import (
	"errors"
	"reflect"
)

func Call(structs interface{}, methodName string, params ...interface{}) (result []reflect.Value, err error) {
	myClassValue := reflect.ValueOf(structs)
	m := myClassValue.MethodByName(methodName)

	if !m.IsValid() {
		return nil, errors.New("method not found param name")
	}

	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	return m.Call(in), nil
}
