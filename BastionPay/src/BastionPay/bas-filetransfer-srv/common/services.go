package common

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
)

type Result struct {
	TotalResult int64       `json:"total_result"`
	HasNext     bool        `json:"has_next"`
	Page        int64       `json:"page"`
	Size        int64       `json:"size"`
	List        interface{} `json:"list"`
}

func (this *Result) PageQuery(
	query *gorm.DB,
	object, list interface{},
	page, size int64,
	class interface{},
	name string) (*Result, error) {

	var total int64

	if size <= 0 {
		size = 100
	}

	offset := (page - 1) * size

	err := query.Model(object).Count(&total).Limit(size).Offset(offset).Find(list).Error
	if err != nil {
		return nil, err
	}

	result := &Result{
		TotalResult: total,
		HasNext:     total > (page * size),
		Page:        page,
		Size:        size,
		List:        list,
	}

	if class != nil && name != "" {
		r, err := this.Calls(class, name, list)
		if err != nil {
			return nil, err
		}

		result.List = r[0].Interface()
	}

	return result, err
}

func (this *Result) PageResult(
	list interface{},
	total,
	page,
	size int64) (*Result, error) {

	return &Result{
		TotalResult: total,
		HasNext:     total > (page * size),
		Page:        page,
		Size:        size,
		List:        list,
	}, nil
}

func (this *Result) Calls(
	myClass interface{},
	name string,
	params ...interface{}) ([]reflect.Value, error) {

	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(name)

	if !m.IsValid() {
		err := fmt.Errorf("method not found param name: %s", name)
		return nil, err
	}

	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	return m.Call(in), nil
}
