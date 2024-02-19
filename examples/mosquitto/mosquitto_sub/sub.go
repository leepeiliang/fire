package main

import (
	mappercommon "fire/pkg/common"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"
	"os"
	"time"
)

var mqttClient mappercommon.MqttClient

func main() {
	var err error

	mqttClient = mappercommon.MqttClient{
		IP:     "tcp://192.168.193.12:1883",
		User:   "root",
		Passwd: "root",
	}
	if err = mqttClient.Connect(); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}

	mqttClient.Subscribe("tongjithree/devices-data-update", onMessage)
	for {
		time.Sleep(time.Second)
	}
}

// onMessage callback function of Mqtt subscribe message.
func onMessage(client mqtt.Client, message mqtt.Message) {
	klog.Info("Receive message: ", message.Topic())

	klog.Info("Device MSg %s: ", string(message.Payload()))

}
