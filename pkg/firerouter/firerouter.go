package firerouter

import (
	"fire/pkg/device"
	"fire/pkg/fire"
	"fire/pkg/fireface"
	"fire/pkg/firenet"
	"fire/pkg/globals"
	"k8s.io/klog/v2"
)

type FireRouter struct {
	firenet.BaseRouter
}

// RequestHandle 根据收到的fire.FireMessage 相应返回。
func (this *FireRouter) RequestHandle(request fireface.IRequest) {
	klog.Infof("Call fireRouter RequestHandle")

	//先读取客户端的数据，再组织返回数据
	//klog.Infof("RequestHandle recv from client :data= %x", request.GetMsgData())

	err := request.GetConnection().SendBuffMsg(request.GetMsgData())
	if err != nil {
		klog.Errorf(err.Error())
	}
}

// DataHandle 根据收到的fire.FireMessage 相应返回。
func (this *FireRouter) DataHandle(request fireface.IRequest) {
	klog.Infof("Call fireRouter DataHandle")
	var data = &fire.FireData{}

	//先读取客户端的数据，再组织返回数据
	klog.Infof("DataHandle recv from client :data= %x", request.GetData())
	if len(request.GetData()) < 8 {
		klog.Infof("垃圾消息无数据位")
		return
	}
	one := globals.GetFireHeart()
	one.HeartProperties(0) //发一个0值心跳
	one.ReSetHeart()       //重新初始化定时心跳
	if len(request.GetData()) == 8 {
		klog.Infof("心跳返回的设备信息是: %x", request.GetData())
		return
	}
	dv := data.DecodeSwitchDataTypeToData(request.GetData())
	//klog.Infof("DecodeSwitchDataTypeToData= %+v", dv)
	//for key, val := range dv {
	//	klog.Infof("DataHandle recv DataValue--KEY:%s Value:%+v", key, val)
	//}
	device.ParseFireToSouth(dv)
}
