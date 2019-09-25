package main

import (
	"BastionPay/merchant-teammanage-api/controllers"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func (this *WebServer) routes() {
	app := this.mIris
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "X-Requested-With", "X_Requested_With", "Content-Type", "Access-Token", "Accept-Language", "Api-Key", "Req-Real-Ip"},
		AllowCredentials: true,
	})

	app.Any("/", func(ctx iris.Context) {
		ctx.JSON(
			map[string]interface{}{
				"code": 0,
			})
	})

	v1 := app.Party("/v1/bk/merchant/teammanage", crs)
	{
		v1.Any("/", func(ctx iris.Context) {
			ctx.JSON(
				map[string]interface{}{
					"code": 0,
				})
		})
		//活动，添加，list，更新，警用 全是后端的
		acParty := v1.Party("/company")
		{
			ac := new(controllers.Company)

			acParty.Post("/add", ac.Add)
			acParty.Post("/get", ac.Get)
			acParty.Post("/del", ac.Del)
			acParty.Post("/update", ac.Update)
			acParty.Post("/list", ac.List)
		}
		//大红包，创建，list（back）
		redParty := v1.Party("/department")
		{
			viplv := new(controllers.Department)

			redParty.Post("/add", viplv.Add)
			redParty.Post("/get", viplv.Get)
			redParty.Post("/del", viplv.Del)
			redParty.Post("/update", viplv.Update)
			redParty.Post("/list", viplv.List)
		}
		//小红包，创建，list（back）
		emParty := v1.Party("/employee")
		{
			viplv := new(controllers.Employee)

			emParty.Post("/add", viplv.Add)
			emParty.Post("/get", viplv.Get)
			emParty.Post("/del", viplv.Del)
			emParty.Post("/update", viplv.Update)
			emParty.Post("/list", viplv.List)
		}

	}

	v1bk := app.Party("/v1/ft/merchant/teammanage", crs)
	{
		v1bk.Any("/", func(ctx iris.Context) {
			ctx.JSON(
				map[string]interface{}{
					"code": 0,
				})
		})

		deParty := v1bk.Party("/department")
		{
			shareInfo := new(controllers.Department)

			deParty.Post("/list", shareInfo.ListFront)
			deParty.Post("/get", shareInfo.GetFront)
			deParty.Post("/gets", shareInfo.GetsFront)
		}
		emParty := v1bk.Party("/employee")
		{
			shareInfo := new(controllers.Employee)

			emParty.Post("/list", shareInfo.ListFront)
			emParty.Post("/get", shareInfo.GetFront)
			emParty.Post("/gets", shareInfo.GetsFront)
		}

	}

}
