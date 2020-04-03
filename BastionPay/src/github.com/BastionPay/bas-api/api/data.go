package api

/////////////////////////////////////////////////////
// input/output data/value
// when input data, user encode and sign data, server decode and verify;
// when output value, server encode and sign data, user decode and verify;
type UserData struct {
	// user unique key
	UserKey string `json:"user_key" doc:"用户唯一标示"`
	// message = origin data -> rsa encode -> base64
	Message string `json:"message" doc:"加密数据，(原始数据->RSA加密)->Base64"`
	// timestamp = Unix timestamp string
	TimeStamp string `json:"time_stamp" doc:"Unix时间戳"`
	// signature = origin data -> rsa encode -> sha512 -> rsa sign -> base64
	Signature string `json:"signature" doc:"签名数据，(原始数据->RSA加密+time_stamp)->sha512->RSA签名->Base64"`
}

// user response/push data
type UserResponseData struct {
	Err    int      `json:"err" doc:"错误码"`     // error code
	ErrMsg string   `json:"errmsg" doc:"错误信息"` // error message
	Value  UserData `json:"value" doc:"返回数据"`  // response data
}
