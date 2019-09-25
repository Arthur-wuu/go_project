package main

import (
	"BastionPay/baspay-encrypte/controllers"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func (this *WebServer) routes() {
	app := this.mIris
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "X-Requested-With", "X_Requested_With", "Content-Type", "Access-Token", "Accept-Language"},
		AllowCredentials: true,
	})

	app.Any("/", func(ctx iris.Context) {
		ctx.JSON(
			map[string]interface{}{
				"code": 0,
			})
	})

	v1 := app.Party("/wallet/api", crs)
	{
		encydata := controllers.EncyData{}
		v1.Any("/{param:path}", encydata.Get)
	}

	v2 := app.Party("/trade", crs)
	{
		form := controllers.Form{}
		v2.Any("/form", form.FormAction)
	}

	//v2 := app.Party("/h5/baspay", crs)
	//{
	//	h5Encydata := controllers.H5EncyData{}
	//	v2.Any("/{param:path}", h5Encydata.H5Get)
	//}

}
