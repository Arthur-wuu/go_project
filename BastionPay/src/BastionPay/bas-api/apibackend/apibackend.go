package apibackend

// router's message format
const (
	DataFromDefault = 0 // from internal
	DataFromApi     = 1 // from /api/ver/srv/function real message
	DataFromUser    = 2 // from /user/ver/srv/function UserMessage(subuserkey+real message)
	DataFromAdmin   = 3 // from /admin/ver/srv/function AdminMessage(subuserkey+real message)

	HttpRouterApi       = "api"
	HttpRouterApiTest   = "apitest"
	HttpRouterUser      = "user"
	HttpRouterUserTest  = "usertest"
	HttpRouterAdmin     = "admin"
	HttpRouterAdminTest = "admintest"
)

// 后台user message的格式
type UserMessage struct {
	SubUserKey string `json:"sub_user_key" doc:"指定用户请求的唯一key"`
	Message    string `json:"message" doc:"实际的请求信息"`
}

// admin message的格式
type AdminMessage struct {
	SubUserKey string `json:"sub_user_key" doc:"指定用户请求的唯一key"`
	Message    string `json:"message" doc:"实际的请求信息"`
}
