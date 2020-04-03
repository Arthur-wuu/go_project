package config

import (
	l4g "github.com/alecthomas/log4go"
)

// srv node config
type ConfigNode struct {
	SrvName    string `json:"srv_name"`    // service name
	SrvVersion string `json:"srv_version"` // service version
	CenterAddr string `json:"center_addr"` // center addr ip:port
}

// load srv node config from absolution path
func (cn *ConfigNode) Load(absPath string) {
	err := LoadJsonNode(absPath, "node", cn)
	if err != nil {
		l4g.Crashf("", err)
	}
}
