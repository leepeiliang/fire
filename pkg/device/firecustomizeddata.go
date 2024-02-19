package device

import (
	"encoding/json"
	"fire/pkg/common"
	"fire/pkg/fire"
	"fire/pkg/globals"
	"k8s.io/klog/v2"
	"strings"
)

const FireAddress = "fireAddress"

func ParseFireToSouth(dvs []fire.Data) {

	var (
		push    bool
		err     error
		topic   = globals.GetHostTopicInfo()
		tmp     = make([]common.DeviceCustomizedData, 0)
		dataMsg common.DeviceGroupCustomizedData
	)
	for id, dv := range dvs {
		klog.Infof("ParseFireToSouth Device:[%d]---len:[%d]", id, len(dv.Data))
		for _, prt := range dv.Data {
			klog.Infof("ParseFireToSouth Device msg:%+v", prt)
		}
		var (
			updateMsg common.DeviceCustomizedData
			tmpDev    *common.BaseDevice
			find      bool
		)
		find = false
		dvtemp := dv

		updateMsg.Timestamp = common.GetTimestamp()
		updateMsg.Data = map[string]*common.DataValue{}

		//klog.V(3).Infof("Address: %s", globals.GetHostReSetCodeInfo())
		for key, val := range dvtemp.Data {
			klog.V(3).Infof("fire.Properties: %s-----%v", key, val)
			if strings.EqualFold(string(key), globals.GetHostReSetCodeInfo()) {
				klog.V(3).Infof("@@@@@@@@@@@@@复位@@@@@@@@@@@@@@@@")
				resetProperties()
				return
			}
			if !find {
				for deviceid, dev := range Devices {
					//klog.V(3).Infof("devices deviceid: %s-----%v", deviceid,dev)
					dataMsg.DeviceType = getDeviceTypeByName(deviceid)
					dataMsg.Timestamp = common.GetTimestamp()

					//klog.V(3).Infof("DData.Properties: %s", key)

					for _, v := range dev.Properties {
						//						klog.V(3).Infof("Dev FireAddress: %s-----%s", key, v.VisitorConfig)
						if strings.EqualFold(string(key), v.VisitorConfig) {
							klog.V(3).Infof("Dev: %s-----%v", deviceid, v.VisitorConfig)
							updateMsg.DeviceID = deviceid
							tmpDev = dev
							find = true
							break
						}
					}

					if find {
						break
					}
				}
			}
			if find {

				for _, v := range tmpDev.Properties {

					if strings.EqualFold(string(key), v.VisitorConfig) {
						klog.V(3).Infof("Dev: Add Property%v", v.VisitorConfig)
						updateMsg.Data[v.PropertyID] = val
						break
					}
				}

			}

		}
		if len(updateMsg.Data) > 0 {
			klog.V(3).Infof("Dev updateMsg: Len(%d)Msg:%v", len(updateMsg.Data), updateMsg)
			tmp = append(tmp, updateMsg)
		}

	}

	if len(tmp) > 0 {
		push = true
		dataMsg.Devices = tmp
	}
	klog.V(3).Infof("Device payload[%v]", dataMsg)
	if push {
		// construct payload
		var payload []byte

		if payload, err = json.Marshal(dataMsg); err != nil {
			klog.Error("Create message data failed")
			return
		}
		//topic := fmt.Sprintf(td.Topic, td.DeviceGroupType)
		klog.V(3).Infof("Device Property[%s]", string(payload))
		if err = globals.MqttClient.Publish(topic, payload); err != nil {
			klog.Error(err)
			return
		}

		klog.V(2).Infof("Update value: %s, topic: %s", dataMsg.DeviceType, topic)
		return
	}

	return
}

// 刷复位

func resetProperties() {
	//klog.Infof("ParseFireToSouth Device:[%d]----%+v", id, dv)
	var (
		push  bool
		err   error
		topic = globals.GetHostTopicInfo()

		dataMsg common.DeviceGroupCustomizedData
	)
	klog.V(3).Infof("devices len: %d", len(Devices))

	for deviceid, dev := range Devices {
		tmpDev := dev
		var tmp = make([]common.DeviceCustomizedData, 0)
		var updateMsg common.DeviceCustomizedData
		updateMsg.Timestamp = common.GetTimestamp()
		updateMsg.Data = map[string]*common.DataValue{}
		klog.V(3).Infof("devices deviceid: %s", deviceid)
		//		klog.V(3).Infof("devices device: %v", dev)
		dataMsg.DeviceType = getDeviceTypeByName(deviceid)
		updateMsg.DeviceID = deviceid
		dataMsg.Timestamp = common.GetTimestamp()

		for _, v := range tmpDev.Properties {
			klog.V(3).Infof("DData.Properties-address: %s", v.VisitorConfig)
			updateMsg.Data[v.PropertyID] = &common.DataValue{
				Value:     0,
				Timestamp: common.GetTimestamp(),
				Metadata: common.DataMetadata{
					Type:      "boolean",
					Timestamp: common.GetTimestamp(),
				},
			}
			if getDataBitKey(v.VisitorConfig, int(fire.PropertyIDKey)) == "0" {
				updateMsg.Data[v.PropertyID] = &common.DataValue{
					Value:     1,
					Timestamp: common.GetTimestamp(),
					Metadata: common.DataMetadata{
						Type:      "boolean",
						Timestamp: common.GetTimestamp(),
					},
				}
			}
			klog.V(3).Infof("DData.Properties-PropertyName: %+v", updateMsg.Data[v.PropertyID])

		}
		if len(updateMsg.Data) > 0 {
			//			klog.V(3).Infof("Dev updateMsg: Len(%d)Msg:%v", len(updateMsg.Data), updateMsg)
			tmp = append(tmp, updateMsg)
		}
		push = false
		if len(tmp) > 0 {
			push = true
			dataMsg.Devices = tmp
		}
		//		klog.V(3).Infof("Device payload[%v]", dataMsg)
		if push {
			// construct payload
			var payload []byte

			if payload, err = json.Marshal(dataMsg); err != nil {
				klog.Error("Create message data failed")
				return
			}
			//topic := fmt.Sprintf(td.Topic, td.DeviceGroupType)
			klog.V(3).Infof("Device Property[%s]", string(payload))
			if err = globals.MqttClient.Publish(topic, payload); err != nil {
				klog.Error(err)
				return
			}

			klog.V(2).Infof("Update value: %s, topic: %s", dataMsg.DeviceType, topic)
		}
	}

}

func getDataBitKey(addrsee string, n int) string {
	tmpSlice := strings.Split(addrsee, "-")
	return tmpSlice[n]
}
func getDeviceTypeByName(name string) string {
	tmpSlice := strings.Split(name, ".")
	return strings.Join(tmpSlice[2:6], ".")
}
func getDeviceIDByName(name string) string {
	tmpSlice := strings.Split(name, ".")
	return strings.Join(tmpSlice[:6], ".")
}
func getPointByDevicePoint(point string) string {
	tmpSlice := strings.Split(point, ".")
	return strings.Join(tmpSlice[6:], ".")
}
