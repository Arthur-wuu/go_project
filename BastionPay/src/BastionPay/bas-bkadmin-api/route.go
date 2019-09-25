package main

import (
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/controllers"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"strings"
)

type (
	Service struct {
		App         *iris.Application
		Config      *tools.Config
		RedisClient *redis.Client
	}
)

func (this *Service) routes() {
	this.App = iris.New()

	this.App.Any("/", new(controllers.Index).Index)

	verify := controllers.Verify{}
	verify.Config = this.Config

	redirectC := controllers.NewRedirectController(this.Config)

	var logCtrl = controllers.NewLogController(this.Config)

	v1 := this.App.Party("/v1", verify.VerifyAccess)
	v1.Done(logCtrl.RecodeLog)

	{
		accountParty := v1.Party("/account")
		{
			accounts := controllers.Account{}

			accountParty.Post("/register", accounts.Register)
			accountParty.Get("/search", accounts.Search)
			accountParty.Get("/user-info", accounts.GetUserInfo)
			accountParty.Get("/batch-user-by-ids", accounts.BatchUserByIds)
			accountParty.Put("/update", accounts.Update)
			accountParty.Put("/disabled", accounts.Disabled)
			accountParty.Put("/set-admin", accounts.SetAdmin)
			accountParty.Put("/before-change-password", accounts.ChangeBeforePassword)
			accountParty.Put("/after-change-password", accounts.ChangeAfterPassword)
			accountParty.Put("/change-user-password", accounts.ChangeUserPassword)
			accountParty.Delete("/delete", accounts.Delete)

			login := controllers.Login{}
			login.Config = this.Config

			accountParty.Post("/login", login.Login)
			accountParty.Delete("/logout", login.Logout)
		}

		accessParty := v1.Party("/access")
		{
			access := controllers.Access{}

			accessParty.Post("/add-access", access.AddAccess)
			accessParty.Get("/search", access.Search)
			accessParty.Delete("/delete", access.Delete)
			accessParty.Put("/update", access.Update)
			accessParty.Get("/search-user-pertain-access", access.SearchUserPertainAccess)
		}

		roleParty := v1.Party("/role")
		{
			role := controllers.Role{}

			roleParty.Post("/add-role", role.AddRule)
			roleParty.Get("/search", role.Search)
			roleParty.Delete("/delete", role.Delete)
			roleParty.Put("/update", role.Update)
			roleParty.Put("/disabled", role.Disabled)
		}

		userRoleParty := v1.Party("/user-role")
		{
			userRole := controllers.UserRole{}

			userRoleParty.Post("/set-user-role", userRole.SetUserRule)
			userRoleParty.Get("/search-user-role", userRole.SearchUserRole)
		}

		roleAccessParty := v1.Party("/role-access")
		{
			roleAccess := controllers.RoleAccess{}

			roleAccessParty.Post("/set-role-access", roleAccess.SetRuleAccess)
			roleAccessParty.Get("/search", roleAccess.Search)
		}

		gaParty := v1.Party("/ga")
		{
			ga := controllers.GA{}
			ga.Config = this.Config

			gaParty.Get("/bind", ga.Bind)
			gaParty.Post("/verify", ga.Verify)
			gaParty.Post("/bind-verify", ga.BindVerify)
		}

		bastionPayParty := v1.Party("/bastionpay")
		{
			bastionPay := controllers.NewBastionPayController(this.Config)

			supportedFunction := make(map[string]string)
			for _, path := range this.Config.WalletPaths {
				index := strings.LastIndex(path, "/")

				relativePath := path[0:index]
				function := path[index+1:]

				supportedFunction[function] = relativePath

				bastionPayParty.Post("/admin"+"/"+function, bastionPay.AdminByFunction)
				fmt.Println("bas=", "/admin"+"/"+function)
			}
		}
		cm := controllers.NewCoinMarket(this.Config)
		pyCm := v1.Party("/coinmarket")
		{ //做代理功能，用Any
			pyCm.Any("/{param:path}", cm.Handler)
		}
		//		pyCm.Post("/ticker", cm.Ticker)

		//rd := controllers.NewRedirectController(this.Config)
		//v1account := v1.Party("/account")
		//{
		//	v1account.Any("/listusers", rd.HandleV1Admin)
		//}
		upFile := controllers.NewUploadFile(this.Config)
		fileParty := v1.Party("/upload")
		{
			fileParty.Any("/coinlogo", upFile.HandleLogoFiles2)
			fileParty.Any("/notice", upFile.HandleNoticeFiles)
			fileParty.Any("/notify", upFile.HandleNotifyFiles)
		}

		//		redirectC := controllers.NewRedirectController(this.Config)
		noticePy := v1.Party("/notice")
		{ //做代理功能，用Any
			noticePy.Any("/{param:path}", redirectC.HandleV1Admin)
		}
		notifyPy := v1.Party("/notify")
		{ //做代理功能，用Any
			notifyPy.Any("/{param:path}", redirectC.HandlerV1BasNotify)
		}
		downLoad := v1.Party("/bkadmin/download")
		{
			downLoadCtrl := controllers.NewDownloadController(this.Config)
			downLoad.Get("/task/status", downLoadCtrl.GetStatus)
		}
		bkadmin := v1.Party("/bkadmin")
		{
			bkadmin.Get("/log/login", logCtrl.GetLoginLog)
			bkadmin.Get("/log/safe", logCtrl.GetOperationLog)
		}
	}

	//v2 := this.App.Party("/v2", verify.VerifyAccess)
	//{
	//	basPy := v2.Party("/bastionpay")
	//	{ //做代理功能，用Any
	//		basPy.Any("/{param:path}", redirectC.HandlerV2Gateway)
	//	}
	//	v2.Any("/account/updatefrozen", redirectC.HandlerV2Gateway)
	//}

	//	test := func(ctx iris.Context) {
	//		//l4g.Info("test CoinMarket recv msg")
	//		ctx.Next()
	//	}
	////	cm := controllers.NewCoinMarket(this.Config)
	//	v9 := this.App.Party("/v1/notify", test)
	//	{ //做代理功能，用Any
	//		v9.Any("/{param:path}", redirectC.HandlerV1BasNotify)
	//	}
	//bastionPay := controllers.NewBastionPayController(this.Config)
	//v92 := this.App.Party("/v1/test", func(ctx iris.Context){ctx.Next()})
	//{ //做代理功能，用Any
	//	v92.Get("/test", bastionPay.Test)
	//}
	//
	//this.App.Post("/v9/ticker", cm.Ticker)
	//
	//rd := controllers.NewRedirectController(this.Config)
	//v9account := this.App.Party("/v1/account", test)
	//{
	//	v9account.Any("/listusers", rd.HandleV1Admin)
	//}
	//upFile := controllers.NewUploadFile(this.Config)
	//fileParty := this.App.Party("/v9/upload", test)
	//{
	//	fileParty.Any("/coinlogo", upFile.HandleLogoFiles2)
	//	fileParty.Any("/notice", upFile.HandleNoticeFiles)
	//}
}
