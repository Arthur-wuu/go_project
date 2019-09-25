// Package classification bastion pay Account API.
//
// Bastion pay user management services include registration, authentication, user information viewing,
// security authentication and other information using JWT token authentication and authentication
// using email, SMS, captcha, Google authentication.
//
//     Schemes: http, https
//     Host: http://account.api.test.mike-huang.cn
//     BasePath: /api/account
//     Version: 0.1.0
//     Contact: Ingram<mike.huang@blockshine.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - JWT:
//
//     SecurityDefinitions:
//     JWT:
//          type: apiKey
//          name: Authorization
//          in: header
//          description: JWT token expiration time is 30 minutes, please re-request before expiration
//
//     Extensions:
//     x-meta-array:
//       - language
//       - timezone
//
// swagger:meta
package main

import (
	"fmt"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-admin-api/controllers"
	"github.com/iris-contrib/middleware/cors"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/i18n"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
)

type App struct {
	Config *config.Config
	Iris   *iris.Application
	Redis  *common.Redis
	Db     *gorm.DB
}

func NewApp(config *config.Config, redis *common.Redis, db *gorm.DB) {
	var (
		a   *App
		app *iris.Application
	)

	app = iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	globalLocale := i18n.New(i18n.Config{
		Default:      "en-US",
		URLParameter: "language",
		Languages: map[string]string{
			"en-US": "./locales/locale_en-US.ini",
			"zh-CN": "./locales/locale_zh-CN.ini"}})
	app.UseGlobal(globalLocale)

	if config.Server.Debug {
		app.Any("/debug/pprof/{action:path}", pprof.New())
	}

	a = &App{
		Config: config,
		Redis:  redis,
		Db:     db,
		Iris:   app,
	}
	a.Controller()
	a.Run()
}

func (a *App) Controller() *App {
	var (
		app              = a.Iris
		interceptorCtrl  = controllers.NewInterceptorController(a.Redis, a.Config)
		verificationCtrl = controllers.NewVerificationController(a.Redis, a.Db, a.Config)
		userCtrl         = controllers.NewUserController(a.Redis, a.Db, a.Config)
		passwordCtrl     = controllers.NewPasswordController(a.Redis, a.Db, a.Config)
		gaCtrl           = controllers.NewGaController(a.Redis, a.Db, a.Config)
		infoCtrl         = controllers.NewInfoController(a.Redis, a.Db, a.Config)
		optionalCtrl     = controllers.NewOptionalController(a.Db)
		logCtrl          = controllers.NewLogController(a.Db, a.Config)
		//keyCtrl          = controllers.NewKeyController(a.Config)
		bastionPayCtrl = controllers.NewBastionPayController(a.Redis, a.Config)
		coinMarketCtrl = controllers.NewCoinMarketController(a.Config)
		//		redirectCtrl   = controllers.NewRedirectController(a.Config)
		noticeCtrl   = controllers.NewNoticeController(a.Db, a.Config)
		downLoadCtrl = controllers.NewDownloadController(a.Config)
	)

	//app.UseGlobal(interceptorCtrl.Interceptor)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "X-Requested-With", "X_Requested_With", "Content-Type", "Access-Token", "Accept-Language"},
		AllowCredentials: true,
	})

	app.Any("/", func(ctx iris.Context) {
		ctx.JSON(map[string]interface{}{
			"code":    0,
			"message": "ok",
			"data":    "",
		})
	})

	party := app.Party("/v1/user/account", crs, interceptorCtrl.Interceptor).AllowMethods(iris.MethodOptions)
	party.Done(logCtrl.RecodeLog)
	//app.DoneGlobal(logCtrl.RecodeLog)

	// 获取验证码
	// swagger:operation GET /verification/{type} verification getVerification
	// ---
	// summary: Get verify code
	// description: Get verify code
	// parameters:
	//   - name: type
	//     in: path
	//     description: type of verification
	//     required: true
	//     type: string
	//     enum:
	//     - email
	//     - sms
	//     - captcha
	//     - ga
	//   - name: operating
	//     in: query
	//     description: type of operating
	//     required: true
	//     type: string
	//     enum:
	//     - login
	//     - register
	//     - forget_password
	//     - withdrawal
	//     - withdrawal_address
	//     - trading
	//     - bind_ga
	//     - unbind_ga
	//     - bind_email
	//     - bind_phone
	//     - rebind_phone
	//   - name: recipient
	//     in: query
	//     description: if type is email,sms and not logged in， need this
	//   - name: captcha_token
	//     in: query
	//     description: if type is email,sms and not logged in, you need to request a captcha token to prevent violent requests
	//     required: false
	//     type: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: string
	party.Get("/verification/{type}", verificationCtrl.Send)

	// 验证验证码
	// swagger:operation POST /verification verification postVerification
	// ---
	// summary: Request to verify
	// description: Request to verify code
	// parameters:
	//   - in: body
	//     name: body
	//     description: require params
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         id:
	//           type: string
	//           format: uuid
	//         value:
	//           type: string
	//           format: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/verification", verificationCtrl.Verification)

	// 刷新token
	// swagger:operation GET /refresh token getRefresh
	// ---
	// summary: Refresh JWT token
	// security:
	// - JWT: []
	// description: Refresh JWT token expiration time before token expires
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//             type: object
	//             properties:
	//               token:
	//                 type: string
	//                 format: byte
	//                 description: JWT token
	//               expiration:
	//                 type: number
	party.Get("/refresh", userCtrl.RefreshToken)
	// 登录
	// swagger:operation POST /login account postLogin
	// ---
	// summary: Login
	// security:
	// - JWT: []
	// description: User username and password to get JWT token
	// parameters:
	//   - in: body
	//     name: body
	//     description: mail and phone need any one
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         username:
	//           type: string
	//           format: string
	//           description: Use the email or phone you used to register
	//         password:
	//           type: string
	//           format: string
	//         captcha_token:
	//           type: string
	//           format: uuid
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: object
	//           properties:
	//             token:
	//               type: string
	//               format: byte
	//               description: JWT token
	//             expiration:
	//               type: number
	//             safe:
	//               type: bool
	//               description: If it is not safe, you need to verify it with /login/ga
	party.Post("/login", userCtrl.Login)

	// 登录
	// swagger:operation POST /login/ga account postLoginWithGa
	// ---
	// summary: Login with ga
	// security:
	// - JWT: []
	// description: Verify ga token
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         ga_token:
	//           type: string
	//           format: uuid
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           properties:
	//             token:
	//               type: string
	//               format: byte
	//               description: JWT token
	//             expiration:
	//               type: number
	party.Post("/login/ga", userCtrl.LoginWithGa)

	// 注册
	// swagger:operation POST /register account postRegister
	// ---
	// summary: Register
	// security:
	// - JWT: []
	// description: Can use email or mobile number to register
	// parameters:
	//   - in: body
	//     name: body
	//     description: email and phone need any one
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         email:
	//           type: string
	//           format: email
	//         phone:
	//           type: string
	//           format: tel
	//         country_code:
	//           type: string
	//           format: string
	//           description: If use phone, you need this
	//         password:
	//           type: string
	//           format: password
	//         citizenship:
	//           type: string
	//           format: string
	//         language:
	//           type: string
	//           format: en-US
	//           description: Get from user's browser
	//         timezone:
	//           type: string
	//           format: +8:00
	//           description: Get from user's browser
	//         token:
	//           type: string
	//           format: uuid
	//           description: Please verify the email or phone
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           properties:
	//             token:
	//               type: string
	//               format: byte
	//               description: JWT token
	//             expiration:
	//               type: number
	party.Post("/register", userCtrl.Register)

	// 用户是否存在
	// swagger:operation POST /exists account postExists
	// ---
	// summary: User exists
	// description: Check if the user exists
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         username:
	//           type: string
	//           format: string
	//         captcha_token:
	//           type: string
	//           format: uuid
	//           description: operating=register
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: bool
	//           format: bool
	party.Post("/exists", userCtrl.Exists)

	// 获取
	// swagger:operation GET /ga ga getGa
	// ---
	// summary: Get google authentication secret
	// security:
	// - JWT: []
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: object
	//           properties:
	//             id:
	//               type: string
	//               format: uuid
	//               description: use post /ga to verify
	//             secret:
	//               type: string
	//               format: string
	//             image:
	//               type: bool
	//               description: Base64 of URI QR code
	party.Get("/ga", gaCtrl.Generate)
	// 绑定
	// swagger:operation POST /ga/bind ga postGa
	// ---
	// summary: Bind google authentication
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         id:
	//           type: string
	//           format: uuid
	//         value:
	//           type: string
	//           format: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/ga/bind", gaCtrl.Bind)
	// 解绑
	// swagger:operation POST /ga/unbind ga unbindGa
	// ---
	// summary: UnBind google authentication
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         email_token:
	//           type: string
	//           format: uuid
	//           description: operating=unbind_ga
	//         sms_token:
	//           type: string
	//           format: uuid
	//           description: operating=unbind_ga
	//         value:
	//           type: string
	//           format: number
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/ga/unbind", gaCtrl.UnBind)

	// 修改
	// swagger:operation POST /password/modify password modifyPassword
	// ---
	// summary: Modify password
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         old_password:
	//           type: string
	//           format: string
	//         password:
	//           type: string
	//           format: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/password/modify", passwordCtrl.Modify)
	// 查询账户
	// swagger:operation POST /password/inquire password inquireResetInfo
	// ---
	// summary: Query reset password information
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         username:
	//           type: string
	//           format: string
	//         captcha_token:
	//           type: string
	//           format: uuid
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: object
	//           properties:
	//             email:
	//               type: string
	//             phone:
	//               type: string
	//             county_code:
	//               type: string
	//             ga:
	//               type: bool
	party.Post("/password/inquire", passwordCtrl.Inquire)
	// 重置
	// swagger:operation POST /password/reset password resetPassword
	// ---
	// summary: Reset password
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         username:
	//           type: string
	//           format: string
	//         password:
	//           type: string
	//           format: string
	//         email_token:
	//           type: string
	//           format: uuid
	//           description: If you have bind email, you need to enter this
	//         sms_token:
	//           type: string
	//           format: uuid
	//           description: If you have bind phone, you need to enter this
	//         ga_value:
	//           type: number
	//           format: int
	//           description: six number like 000000
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/password/reset", passwordCtrl.Reset)
	// 读取信息
	// swagger:operation GET /info information getInfo
	// ---
	// summary: Account information
	// security:
	// - JWT: []
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Get("/info", infoCtrl.GetInformation)
	party.Get("/info/nohide", infoCtrl.GetInformationNoHide)

	// 修改信息
	// swagger:operation POST /info information updateInfo
	// ---
	// summary: Update information
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         language:
	//           type: string
	//         timezone:
	//           type: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/info", infoCtrl.SetInformation)
	// 绑定邮箱
	// swagger:operation POST /info/email information bindEmail
	// ---
	// summary: Bind email
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         email_token:
	//           type: string
	//           format: uuid
	//         email:
	//           type: string
	//           format: email
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/info/email", infoCtrl.BindEmail)
	// 绑定手机
	// swagger:operation POST /info/phone information bindPhone
	// ---
	// summary: Bind phone
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         sms_token:
	//           type: string
	//           format: uuid
	//         phone:
	//           type: string
	//         country_code:
	//           type: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/info/phone", infoCtrl.BindPhone)

	// 绑定手机
	// swagger:operation POST /info/phone/rebind information RebindPhone
	// ---
	// summary: Rebind phone
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         email_token:
	//           type: string
	//           format: uuid
	//           description: operating=rebind_phone
	//         ga_token:
	//           type: string
	//           format: uuid
	//           description: operating=rebind_phone
	//         old_sms_token:
	//           type: string
	//           format: uuid
	//           description: operating=rebind_phone
	//         new_sms_token:
	//           type: string
	//           format: uuid
	//           description: operating=bind_phone
	//         phone:
	//           type: string
	//         country_code:
	//           type: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/info/phone/rebind", infoCtrl.RebindPhone)

	// 读取自选
	// swagger:operation GET /optional optional getOptional
	// ---
	// summary: Get optional
	// security:
	// - JWT: []
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: array
	//           items:
	//             type: string
	party.Get("/optional", optionalCtrl.Get)
	// 修改自选
	// swagger:operation POST /optional optional updateOptional
	// ---
	// summary: Update optional
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         value:
	//           type: array
	//           items:
	//             type: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/optional", optionalCtrl.Update)
	// 登录日志
	// swagger:operation GET /log/login log getLoginLog
	// ---
	// summary: Login log
	// security:
	// - JWT: []
	// parameters:
	//   - name: page
	//     in: query
	//     description: Which page, start from 1
	//     default: 1
	//     required: false
	//     type: number
	//   - name: limit
	//     in: query
	//     description: One page record number, between 1 and 100
	//     default: 10
	//     required: false
	//     type: number
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: object
	//           properties:
	//             total:
	//               type: number
	//               description: Total number
	//             page:
	//               type: number
	//               description: Current page
	//             data:
	//               type: array
	party.Get("/log/login", logCtrl.GetLoginLog)

	// 安全操作日志
	// swagger:operation GET /log/safe log getSafeLog
	// ---
	// summary: Security settings log
	// security:
	// - JWT: []
	// parameters:
	//   - name: page
	//     in: query
	//     description: Which page, start from 1
	//     default: 1
	//     required: false
	//     type: number
	//   - name: limit
	//     in: query
	//     description: One page record number, between 1 and 100
	//     default: 10
	//     required: false
	//     type: number
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: object
	//           properties:
	//             total:
	//               type: number
	//               description: Total number
	//             page:
	//               type: number
	//               description: Current page
	//             data:
	//               type: array
	party.Get("/log/safe", logCtrl.GetOperationLog)

	// 上传公钥--废弃，统一使用bastionpay/user接口
	// swagger:operation POST /report/key optional report key
	// ---
	// summary: report key
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         public_key:
	//           type: string
	//         source_ip:
	//           type: string
	//         callback_url:
	//           type: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	//party.Post("/report/key", keyCtrl.Report)

	// 查询bastionpay
	// swagger:operation POST /bastionpay/query query bastionpay
	// ---
	// summary: query bastionpay
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         function:
	//           type: string
	//         message:
	//           type: string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           type: null
	party.Post("/bastionpay/user", bastionPayCtrl.User)
	//next := func(ctx iris.Context) {
	//	fmt.Println("test recv msg")
	//	ctx.Next()
	//}

	v1 := app.Party("/v1/user", crs, interceptorCtrl.Interceptor).AllowMethods(iris.MethodOptions)
	bp := v1.Party("/bastionpay")
	{
		bp.Any("/{param:path}", bastionPayCtrl.HandlerV1)
	}
	downLoad := v1.Party("/download")
	{
		downLoad.Get("/task/status", downLoadCtrl.GetStatus)
	}
	//pushParty := v1.Party("/push")
	//{
	//	pushParty.Any("/{param:path}", redirectCtrl.HandlerV1Gateway)
	//}

	// 查询币种行情
	// swagger:operation Post /coinmarket/ticker query coinmarket
	// ---
	// summary: Query coinmarket quote
	// security:
	// - JWT: []
	// parameters:
	//   - in: body
	//     name: body
	//     description:
	//     required: true
	//     schema:
	//       type: object
	//       properties:
	//         coinsymbols:
	//           type: []string
	//         converts:
	//           type: []string
	// responses:
	//   "200":
	//     description: success response
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: object
	//           properties:
	//             code:
	//               type: number
	//             msg:
	//               type: string
	//         result:
	//           data: []string
	party.Post("/coinmarket/ticker", coinMarketCtrl.Ticker)

	//内部接口，不需要权限
	//cmPartyNoAuth := app.Party("/v1/coinmarket", crs, func(ctx iris.Context) { ctx.Next() }).AllowMethods(iris.MethodOptions)
	//{
	//	cmPartyNoAuth.Any("/ticker", coinMarketCtrl.Ticker)
	//}

	partyNoAuth := app.Party("/v1/inner/user/account", crs, func(ctx iris.Context) { ctx.Next() }).AllowMethods(iris.MethodOptions)
	{
		partyNoAuth.Get("/getuserinfo", infoCtrl.GetUserInfo)
		partyNoAuth.Post("/listusers", infoCtrl.Listusers)
	}
	//
	//tt := func(ctx iris.Context) {
	//	appClaims := &common.AppClaims{}
	//	appClaims.UserId = 50051
	//	appClaims.Safe = true
	//	appClaims.Uuid = "d06327ca-0d05-42a2-898c-4833e6238d20"
	//	ctx.Values().Set("app_claims", appClaims)
	//	fmt.Println("get one")
	//	ctx.Next()
	//}

	noticePartyNoAuth := app.Party("/v1/inner/user/notice", crs, func(ctx iris.Context) { ctx.Next() }).AllowMethods(iris.MethodOptions)
	{
		noticePartyNoAuth.Post("/getlist", noticeCtrl.GetListFromInner)
		noticePartyNoAuth.Post("/get", noticeCtrl.GetFromInner)
		noticePartyNoAuth.Post("/add", noticeCtrl.AddFromInner)
		noticePartyNoAuth.Post("/update", noticeCtrl.UpdateFromInner)
		noticePartyNoAuth.Post("/del", noticeCtrl.DelFromInner)
		noticePartyNoAuth.Post("/focus", noticeCtrl.FocusFromInner)
		//noticePartyNoAuth.Post("/getlist", noticeCtrl.GetList)
		//noticePartyNoAuth.Post("/get", noticeCtrl.Get)
		//noticePartyNoAuth.Post("/count", noticeCtrl.CountUserNotices)
	}

	noticeParty := app.Party("/v1/user/notice", crs, interceptorCtrl.Interceptor).AllowMethods(iris.MethodOptions)
	{
		noticeParty.Post("/getlist", noticeCtrl.GetList)
		noticeParty.Post("/get", noticeCtrl.Get)
		noticeParty.Post("/count", noticeCtrl.CountUserNotices)
	}
	//testParty := app.Party("/testinfo", crs, tt).AllowMethods(iris.MethodOptions)
	//{
	//	testParty.Get("/info", infoCtrl.GetInformation)
	//}
	//
	//// 获取
	//app.Get("/apikey")
	//// 生成
	//app.Post("/apikey")
	//// 禁用
	//app.Delete("/apikey")
	//
	//// 获取验证码

	/**
	登录
		验证输入
		验证密码
		验证GA
		生成并返回token
		记录登录日志

	注册
		验证输入
		验证邮箱
		创建密码
		创建账户
		生成并返回token
		记录登录日志

	读取基本信息
		验证token
		获取并返回信息

	修改密码
		验证token
		验证输入
		验证密码
		验证邮箱或GA
		创建密码
		修改账户

	找回密码
		验证输入
		验证邮箱
		验证GA
		创建密码
		修改账户

	获取GA绑定URI
		验证token
		生成并存储URI
		返回URI

	绑定GA(Google authenticator)
		验证token
		验证输入
		验证邮箱
		验证密码
		验证GA
		绑定GA

	解绑GA
		验证token
		验证输入
		验证邮箱
		验证密码
		验证GA
		解绑GA

	生成apikey
		验证token
		验证邮箱
		验证GA
		生成并存储Apikey
		返回Apikey

	读取apikey
		验证token
		读取并返回Apikey

	禁用apikey
		验证token
		验证GA
		验证邮箱
		禁用Apikey

	获取验证码
	验证验证码
	*/
	return a
}

func (a *App) Run() {
	fmt.Printf("Server version: %s\n", Version)
	a.Iris.Run(iris.Addr(":" + a.Config.Server.Port))
}
