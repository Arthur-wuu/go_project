package v1

// 冻结
type ReqCpFrozen struct {
	Reason string `json:"reason" doc:"冻结原因"`
}

type AckCpFrozen struct {
	IsFrozen int `json:"is_frozen" doc:"冻结状态"`
}