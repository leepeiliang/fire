package synccomfig

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fire/config"
	"fire/pkg/client"
	mappercommon "fire/pkg/common"
	"fire/pkg/device"
	"fire/pkg/globals"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"time"
)

const Fire = "fire"

const (
	host   = "http://192.168.233.32:30090"
	Active = "/v1/driver/sync"
)

// SyncConfig
type SyncConfig interface {
	SyncConfigActive(ctx context.Context) (*SyncResponse, error)
}

type syncConfig struct {
	client http.Client
}

// NewClient 初始化链接对象
func NewClient(conf client.Config) SyncConfig {
	return &syncConfig{
		client: client.New(conf),
	}
}

type SyncRequest struct {
	Node     string `json:"hostname" description:"节点"`
	Protocol string `json:"protocol" description:"协议"`
	Status   string `json:"status" description:"状态"`
	Pages    int    `json:"pages"`
	Limit    int    `json:"limit"`
}

type SyncResponse struct {
	Total   int64                              `json:"total" description:"设备总数"`
	Devices map[string]mappercommon.BaseDevice `json:"devices" description:"设备"`
}

func (aha *syncConfig) SyncConfigActive(ctx context.Context) (*SyncResponse, error) {
	hostname := globals.GetHostNameInfo()
	if hostname == "" {
		return nil, errors.New("get system hostname err")
	}
	klog.V(2).Infof("System:hostname:%v", hostname)
	roleReq := &SyncRequest{
		Node:     hostname,
		Protocol: Fire,
		Pages:    1,
		Limit:    5000,
	}

	devices := &SyncResponse{}
	//var devices = make(map[string]*mappercommon.BaseDevice)
	address := fmt.Sprintf("http://%s:%d", config.DefaultConfig.EdgeServer.Host, config.DefaultConfig.EdgeServer.Port)
	err := client.POST(ctx, &aha.client, address+Active, roleReq, devices)
	if err != nil {
		klog.Errorf("POST:%s ", err.Error())
		return nil, err
	}
	klog.V(4).Infof("ResponseResult:back:%+v", devices)
	return devices, nil
}

var (
	ctx = context.Background()
)

func FirstSyncConfig() {
	cli := NewClient(client.Config{
		Timeout:      30,
		MaxIdleConns: 100,
	})

	result, err := cli.SyncConfigActive(ctx)
	if err != nil {
		klog.Errorf("SyncConfig%+v", err.Error())
		return
	}

	klog.V(2).Infof("Devices:len:%d", len(result.Devices))
	klog.V(2).Infof("Devices:Total:%d", result.Total)
	klog.V(2).Infof("configmap path:%s", config.DefaultConfig.Configmap)

	payload, err := json.Marshal(result.Devices)
	if err != nil {
		klog.Errorf("Devices:Marshal:%s", err.Error())
		return
	}

	klog.V(2).Infof("Device data len: %d", len(payload))
	//	klog.V(4).Infof("Device data: %s", string(payload))
	err = ioutil.WriteFile(config.DefaultConfig.Configmap, payload, 0666)
	if err != nil {
		klog.Errorf("Devices WriteFile:%s ", err.Error())
		return
	}
	jsonFile, err := ioutil.ReadFile(config.DefaultConfig.Configmap)
	if err != nil {
		klog.Errorf("Devices ReadFile:%s ", err.Error())
		return
	}
	klog.V(2).Info("Devices-ReadFile-len:%d", len(jsonFile))
	jsonFile = bytes.TrimPrefix(jsonFile, []byte("\xef\xbb\xbf"))
	for deviceId, _ := range device.Devices {
		device.Devices[deviceId] = nil
		delete(device.Devices, deviceId) //将小明:100从map中删除
	}
	if err = json.Unmarshal(jsonFile, &device.Devices); err != nil {
		klog.Errorf("Devices Unmarshal:%s ", err.Error())
		return

	}
	klog.V(2).Infof("Upate DeviceInstances Success :len:[%d]", len(device.Devices))

}

func SubscribeSyncConfigMap() {

	FirstSyncConfig()

	hostname := "softwarepark"
	//hostname := globals.GetHostNameInfo()
	//if hostname == "" {
	//	klog.Errorf("POST:%s ", errors.New("get system hostnam err"))
	//	return
	//}
	klog.Infof("System:hostname:%v", hostname)
	topic := fmt.Sprintf("%s/devices-data-update", hostname)
	err := globals.MqttClient.Subscribe(topic, onMessage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Subscribe topic success")
	}
}

type SyncConifg struct {
	Protocol    string    `json:"protocol"`
	PublishTime time.Time `json:"publish_time"`
}

// onMessage callback function of Mqtt subscribe message.
func onMessage(client mqtt.Client, message mqtt.Message) {
	klog.Info("Receive message", message.Topic())
	klog.Info("Device MSg %s: ", string(message.Payload()))

	var sync SyncConifg
	err := json.Unmarshal(message.Payload(), &sync)
	if err != nil {
		klog.Errorf("Unmarshal:%s ", errors.New("Unmarshal sync protocol "))
		return
	}
	if sync.Protocol == "fire" && sync.PublishTime.Add(30*time.Second).After(time.Now().UTC()) {
		time.Sleep(25 * time.Second)
		FirstSyncConfig()
	}
}
