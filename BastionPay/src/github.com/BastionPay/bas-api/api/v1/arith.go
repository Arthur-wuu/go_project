package v1

type Args struct {
	A int `json:"a" doc:"加数1"`
	B int `json:"b" doc:"加数2"`
}

type AckArgs struct {
	C int `json:"c" doc:"和"`
}
