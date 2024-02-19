package fire

import (
	"k8s.io/klog/v2"
	"testing"
)

func TestDataUnmarshal(t *testing.T) {
	var login = []byte{
		0x08, 0x01, 0x00, 0x00, 0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16,
	}

	var fireMessage = &FireData{}
	fireMessage.Decode(login)
	klog.Infof("FireData:%v", fireMessage)

	fireMessage.FireBuildFacilitiesTimeStatDecodeToData()
}
