package device

import (
	"BastionPay/merchant-api/config"
	"strings"
)

var GDeviceMgr DeviceMgr

type DeviceMgr struct {
	mDeviceMap map[string]Device
}

func (this *DeviceMgr) Init() error {
	this.mDeviceMap = make(map[string]Device)
	for i := 0; i < len(config.GConfig.Devices); i++ {

		devconf := config.GConfig.Devices[i]
		switch strings.ToLower(devconf.Name) {
		case "game":
			game := new(Game)
			if err := game.Init(devconf.Addr, devconf.Id); err != nil {
				//ws 设备可能不在线，init开启ws的start
				continue
			}
			this.Add(devconf.Id, game)
		case "doll":
			game := new(Game)
			if err := game.Init(devconf.Addr, devconf.Id); err != nil {
				continue
			}
			this.Add(devconf.Id, game)
		default:

		}
	}
	return nil
}

func (this *DeviceMgr) Add(id string, d Device) {
	this.mDeviceMap[id] = d
}

func (this *DeviceMgr) Get(id string) Device {
	return this.mDeviceMap[id]
}

type Device interface {
	Init(addr, id string) error
	GetId() string
	Send(data interface{}) error
}
