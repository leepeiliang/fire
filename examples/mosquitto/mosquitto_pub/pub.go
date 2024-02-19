package main

import (
	"encoding/json"
	mappercommon "fire/pkg/common"
	"k8s.io/klog/v2"
	"os"
	"sync"
)

var mqttClient mappercommon.MqttClient
var wg sync.WaitGroup

func main() {
	var err error
	push := true

	mqttClient = mappercommon.MqttClient{
		IP:     "tcp://192.168.193.12:1883",
		User:   "root",
		Passwd: "root",
	}
	if err := mqttClient.Connect(); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}

	var (
		dataMsg   mappercommon.DeviceGroupCustomizedData
		updateMsg mappercommon.DeviceCustomizedData
		tmp       = make([]mappercommon.DeviceCustomizedData, 0)
	)
	updateMsg.Timestamp = mappercommon.GetTimestamp()
	updateMsg.Data = map[string]*mappercommon.DataValue{}

	updateMsg.DeviceID = "34.0.1.4.8.4"
	dataMsg.Timestamp = mappercommon.GetTimestamp()
	updateMsg.Data["34.0.1.4.8.4.2.13.1.1"] = &mappercommon.DataValue{
		Value:     1,
		Timestamp: mappercommon.GetTimestamp(),
		Metadata: mappercommon.DataMetadata{
			Type:      "boolean",
			Timestamp: mappercommon.GetTimestamp(),
		},
	}
	if len(updateMsg.Data) > 0 {
		tmp = append(tmp, updateMsg)
	}
	push = false
	if len(tmp) > 0 {
		push = true
		dataMsg.Devices = tmp
	}
	if push {
		// construct payload
		var payload []byte

		if payload, err = json.Marshal(dataMsg); err != nil {
			klog.Error("Create message data failed")
			return
		}
		//topic := fmt.Sprintf(td.Topic, td.DeviceGroupType)
		klog.Infof("Device pproperty[%s]", string(payload))
		if err = mqttClient.Publish("tongjithree/devices-data-update", payload); err != nil {
			klog.Error(err)
			return
		}

		klog.Infof("Update topic: %s", "tongjithree/devices-data-update")
	}
}
