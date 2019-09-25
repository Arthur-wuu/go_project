package main

import (
	"BastionPay/bas-filetransfer-srv/controllers"
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

	v1 := app.Party("/v1", crs)
	{
		bulletin := v1.Party("/filetransfer")
		{

			bulletinController := controllers.Bulletin{}

			bulletin.Post("/export", bulletinController.Export)
			bulletin.Get("/cancel", bulletinController.Cancel)
			bulletin.Get("/status/get", bulletinController.GetStatus)
			bulletin.Post("/status/add", bulletinController.AddStatus)
			bulletin.Post("/status/update", bulletinController.UpdateStatus)
		}
	}

}
