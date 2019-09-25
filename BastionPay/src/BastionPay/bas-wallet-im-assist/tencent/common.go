package tencent

type (
	Request struct {
		ToAccount []string `json:"To_Account,omitempty"`
	}

	Response struct {
		ActionStatus string  `json:"ActionStatus,omitempty"`
		ErrorCode    int     `json:"ErrorCode,omitempty"`
		ErrorInfo    string  `json:"ErrorInfo,omitempty"`
		QueryResult  []Query `json:"QueryResult,omitempty"`
	}

	Query struct {
		ToAccount string `json:"To_Account,omitempty"`
		State     string `json:"State,omitempty"`
	}
)

//往tencent查询用户状态  在task里实现了
func (this *Request) Send() (*Response, error) {
	return &Response{}, nil
}
