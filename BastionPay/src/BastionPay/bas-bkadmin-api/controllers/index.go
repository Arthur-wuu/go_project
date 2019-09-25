package controllers

import "github.com/kataras/iris"

type Index struct {
	Controllers
}

func (this *Index) Index(ctx iris.Context) {
	ctx.JSON(Response{Message: "iris requests server ok"})
}
