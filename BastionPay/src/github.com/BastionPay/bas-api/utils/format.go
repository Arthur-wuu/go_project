package utils

import (
	"encoding/json"
)

func FormatSample(data interface{}) string {
	var (
		result string
	)

	if data == nil{
		result = ""
	} else if d, ok := data.([]byte); ok {
		result = string(d)
	} else if d, ok := data.(string); ok {
		result = d
	} else{
		example, _ := json.Marshal(data)
		result = string(example)
	}

	return result
}

func FormatComment(data interface{}) string {
	return FieldTag(data, 0)
}
