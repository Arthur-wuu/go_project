package common

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
)

func (t *Tools) ShouldBindJSON(ctx iris.Context, params interface{}) error {
	err := ctx.ReadJSON(params)
	if err != nil {
		return err
	}

	ok, err := govalidator.ValidateStruct(params)
	if !ok || err != nil {
		return err
	}

	return nil
}

func (t *Tools) ShouldBindQuery(ctx iris.Context, params interface{}) error {
	err := ctx.ReadForm(params)
	if err != nil {
		return err
	}

	ok, err := govalidator.ValidateStruct(params)
	if !ok || err != nil {
		return err
	}

	return nil
}
