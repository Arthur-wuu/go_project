package common

import (
	"errors"
	"github.com/kataras/iris"
	"reflect"
)

type PageParams struct {
	Page  int `params:"page"`
	Limit int `params:"limit"`
}

type TimeIntervalParams struct {
	StartTime int64 `params:"start_time"`
	EndTime   int64 `params:"end_time"`
}

type ParamsVerification struct {
	Condition bool
	ErrMsg    string
}

func CheckParams(pv *[]ParamsVerification) *ParamsVerification {
	for _, p := range *pv {
		if p.Condition {
			return &p
		}
	}

	return nil
}

func GetParams(ctx iris.Context, params interface{}) error {
	var (
		elem reflect.Value
		err  error
	)

	elem = reflect.ValueOf(params).Elem()
	err = reflectAction(ctx, elem)
	if err != nil {
		return err
	}

	return nil
}

func reflectAction(ctx iris.Context, value reflect.Value) error {
	var (
		err   error
		count int
		tag   string
		to    reflect.Type
	)

	to = value.Type()
	count = value.NumField()

	for i := 0; i < count; i++ {
		f := value.Field(i)

		if f.IsValid() && f.CanSet() {
			tag = string(to.Field(i).Tag.Get("params"))

			switch f.Kind() {
			case reflect.String:
				f.SetString(ctx.URLParam(tag))
			case reflect.Uint:
				pa, err := ctx.URLParamInt(tag)
				if err != nil {
					break
				}
				f.Set(reflect.ValueOf(uint(pa)))
			case reflect.Uint64:
				pa, err := ctx.URLParamInt(tag)
				if err != nil {
					break
				}
				f.Set(reflect.ValueOf(uint64(pa)))
			case reflect.Int:
				pa, err := ctx.URLParamInt(tag)
				if err != nil {
					break
				}
				f.Set(reflect.ValueOf(pa))
			case reflect.Int64:
				pa, err := ctx.URLParamInt64(tag)
				if err != nil {
					break
				}
				f.SetInt(pa)
			case reflect.Float64:
				pa, err := ctx.URLParamFloat64(tag)
				if err != nil {
					break
				}
				f.SetFloat(pa)
			case reflect.Bool:
				pa, err := ctx.URLParamBool(tag)
				if err != nil {
					break
				}
				f.SetBool(pa)
			case reflect.Ptr:
				ctx.URLParams()
				err = reflectAction(ctx, f)
				if err != nil {
					return err
				}
			}

			if err != nil {
				return errors.New("Parsing" + tag + "failed")
			}
		}
	}

	return nil
}
