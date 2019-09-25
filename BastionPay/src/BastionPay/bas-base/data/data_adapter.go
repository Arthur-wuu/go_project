package data

import (
	"BastionPay/bas-api/api"
	"strings"
)

// method
func (method *SrvMethod) FromPath(path string) {
	path = strings.TrimLeft(path, "/")
	path = strings.TrimRight(path, "/")
	paths := strings.Split(path, "/")
	for i := 0; i < len(paths); i++ {
		if i == 1 {
			method.Version = paths[i]
		} else if i == 2 {
			method.Srv = paths[i]
		} else if i >= 3 {
			if method.Function != "" {
				method.Function += "."
			}
			method.Function += paths[i]
		}
	}
}

// data
func (data *SrvData) FromApiData(ud *api.UserData) {
	data.UserKey = ud.UserKey
	data.SubUserKey = ""
	data.Message = ud.Message
	data.TimeStamp = ud.TimeStamp
	data.Signature = ud.Signature
}

func (data *SrvData) ToApiData(ud *api.UserData) {
	ud.UserKey = data.UserKey
	ud.Message = data.Message
	ud.TimeStamp = data.TimeStamp
	ud.Signature = data.Signature
}

// response
func (response *SrvResponse) ToApiResponse(ur *api.UserResponseData) {
	ur.Err = response.Err
	ur.ErrMsg = response.ErrMsg
	response.Value.ToApiData(&ur.Value)
}
