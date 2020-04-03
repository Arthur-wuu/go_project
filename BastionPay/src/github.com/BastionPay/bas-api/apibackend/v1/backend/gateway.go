package backend

// service info
type ServiceInfo struct {
	Version string `json:"version" doc:"服务版本"` // srv version
	Srv     string `json:"srv" doc:"服务名称"`     // srv name
	Count   int    `json:"count" doc:"服务个数"`   // srv 个数
}

// service info
type ServiceInfoList struct {
	Data []ServiceInfo `json:"data"" doc:"服务列表"`
}
