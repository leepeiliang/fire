package main

import (
	"fire/pkg/device"
	"fire/pkg/synccomfig"
	"flag"
	"fmt"
	"os"
	"runtime"

	"fire/config"
	"fire/pkg/common"
	"fire/pkg/fire"
	"fire/pkg/fireface"
	"fire/pkg/firenet"
	"fire/pkg/firerouter"
	"fire/pkg/globals"

	"k8s.io/klog/v2"
)

// DoConnectionBegin 创建连接的时候执行
func DoConnectionBegin(conn fireface.IConnection) {
	klog.Infof("DoConnecionBegin is Called ... ")
	//
	////设置两个链接属性，在连接创建之后
	//klog.Infof("Set conn Name, Home done!")
	//conn.SetProperty("Name", "Aceld")
	//conn.SetProperty("Home", "https://www.kancloud.cn/@aceld")
	//
	//err := conn.SendMsg([]byte("DoConnection BEGIN..."))
	//if err != nil {
	//	klog.Error(err)
	//}
}

// DoConnectionLost 连接断开的时候执行
func DoConnectionLost(conn fireface.IConnection) {
	////在连接销毁之前，查询conn的Name，Home属性
	//if name, err := conn.GetProperty("Name"); err == nil {
	//	klog.Error("Conn Property Name = ", name)
	//}
	//
	//if home, err := conn.GetProperty("Home"); err == nil {
	//	klog.Error("Conn Property Home = ", home)
	//}

	klog.Infof("DoConneciotnLost is Called ... ")
}

func main() {
	var (
		fs1 flag.FlagSet
		err error
	)

	runtime.GOMAXPROCS(2)

	klog.InitFlags(nil)
	fs1.Set("log_dir", "test.md/")
	fs1.Set("log_file", "fire.logs")
	fs1.Set("add_dir_header", "true")
	fs1.Set("logtostderr", "false")
	defer klog.Flush()

	fs1.PrintDefaults()
	if err := config.DefaultConfig.Parse(); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
	klog.V(0).Info(config.DefaultConfig.Configmap)

	globals.MqttSubscribeClient = common.MqttClient{
		IP:         config.DefaultConfig.Mqtt.ServerAddress,
		User:       config.DefaultConfig.Mqtt.UserName,
		Passwd:     config.DefaultConfig.Mqtt.Password,
		Cert:       config.DefaultConfig.Mqtt.CertFile,
		PrivateKey: config.DefaultConfig.Mqtt.PrivateKeyFile,
		Qos:        byte(config.DefaultConfig.Mqtt.Qos),
		Retained:   config.DefaultConfig.Mqtt.Retained}
	globals.MqttPublishClient = globals.MqttSubscribeClient

	subscribeClientId := fmt.Sprintf("bacnet-mapper-subscribe-%s", config.DefaultConfig.Mqtt.UserName)

	if err = globals.MqttSubscribeClient.Connect(subscribeClientId); err != nil {
		klog.Errorf("PrepareWork: mqtt connect err: %s\n", err)
		os.Exit(1)
	}
	publishClientId := fmt.Sprintf("bacnet-mapper-publish-%s", config.DefaultConfig.Mqtt.UserName)

	if err = globals.MqttPublishClient.Connect(publishClientId); err != nil {
		klog.Errorf("PrepareWork: mqtt connect err: %s\n", err)
		os.Exit(1)
	}

	klog.V(0).Info("启动消防HTTP客户端,获取初始化参数")
	go synccomfig.SubscribeSyncConfigMap()
	//klog.V(0).Info("启动消防HTTP服务")
	//go firehttp.FireServer()

	//// 启动定时器
	klog.V(0).Info("初初始化消防定时服务")
	globals.New()
	if err := device.DevInit(config.DefaultConfig.Configmap); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
	//device.DevStart()
	//创建一个server句柄
	s := firenet.NewServer(config.DefaultConfig.Server)
	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由
	klog.V(0).Info("启动消防TCP服务")
	s.AddRouter(fire.DefaultFireMsgID, &firerouter.FireRouter{})
	//开启服务
	s.Serve()

	return
}
