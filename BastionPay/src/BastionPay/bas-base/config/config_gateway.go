package config

import (
	l4g "github.com/alecthomas/log4go"
)

// api gateway config
type ConfigGateway struct {
	Port           string `json:"port"`            // http port
	GatewayVersion string `json:"gateway_version"` // gateway version
	GatewayName    string `json:"gateway_name"`    // gateway name
	GatewayPort    string `json:"gateway_port"`    // gateway rpc port
	TestMode       int    `json:"test_mode"`
	EnableRemoteIp bool   `json:"enable_remote_ip"`
}

// load gateway config from absolution path
func (cc *ConfigGateway) Load(absPath string) {
	err := LoadJsonNode(absPath, "gateway", cc)
	if err != nil {
		l4g.Crashf("", err)
	}
}
