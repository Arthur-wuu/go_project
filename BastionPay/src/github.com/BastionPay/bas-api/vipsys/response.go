package vipsys


type(
	Response struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}

	ResponseLevel struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    *Level      `json:"data,omitempty"`
	}

	ResponseLevels struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    []*Level     `json:"data,omitempty"`
	}

	ResponseUserLevel struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    *UserLevel      `json:"data,omitempty"`
	}

	ResponseUserLevels struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    []*UserLevel     `json:"data,omitempty"`
	}
	ResponseUserRule struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    *ResultUserRule      `json:"data,omitempty"`
	}

	ResponseUserRules struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    []*ResultUserRule     `json:"data,omitempty"`
	}
)
