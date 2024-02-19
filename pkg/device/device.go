package device

import (
	mappercommon "fire/pkg/common"
	"fire/pkg/configmap"
	"sync"
)

const (
	Interface  = "interface"
	Port       = "port"
	MaxPDU     = "maxPDU"
	MinDevice  = "minDevice"
	MaxDevice  = "maxDevice"
	DeviceType = "deviceType"
	Timestamp  = "timestamp"
)

var Devices = make(map[string]*mappercommon.BaseDevice)
var wg sync.WaitGroup

// DevInit initialize the device datas.
func DevInit(configmapPath string) error {
	return configmap.NewParse(configmapPath, Devices)
}
