package firehttp

import (
	"bytes"
	"encoding/json"
	"fire/config"
	mappercommon "fire/pkg/common"
	"fire/pkg/device"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
)

func FireServer() {

	// 1.创建路由
	// 默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()
	// JSON绑定
	r.POST("v1/configmap/fire", func(c *gin.Context) {
		// 声明接收的变量
		var devices = make(map[string]*mappercommon.BaseDevice)
		// 将request的body中的数据，自动按照json格式解析到结构体
		if err := c.ShouldBindJSON(&devices); err != nil {
			// 返回错误信息
			// gin.H封装了生成json数据的工具

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//klog.Infof("Devices:%v",device.Devices)
		klog.Infof("Devices:len:%d", len(devices))

		klog.Infof("configmap path:%s", config.DefaultConfig.Configmap)

		payload, err := json.Marshal(devices)
		if err != nil {
			// 返回错误信息
			// gin.H封装了生成json数据的工具
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		klog.V(2).Infof("Device data len: %d", len(payload))
		klog.V(2).Infof("Device data: %s", string(payload))
		err = ioutil.WriteFile(config.DefaultConfig.Configmap, payload, 0666)
		if err != nil {
			klog.Errorf("configmap Parse:%s ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		jsonFile, err := ioutil.ReadFile(config.DefaultConfig.Configmap)
		if err != nil {
			klog.Errorf("configmap Parse:%s ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		klog.Info("body-len:%d", len(jsonFile))
		jsonFile = bytes.TrimPrefix(jsonFile, []byte("\xef\xbb\xbf"))
		for deviceId, _ := range device.Devices {
			device.Devices[deviceId] = nil
			delete(device.Devices, deviceId) //将小明:100从map中删除
		}
		if err = json.Unmarshal(jsonFile, &device.Devices); err != nil {
			klog.Errorf("%v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		klog.Infof("Upate DeviceInstances Success :len:[%d]", len(device.Devices))

		// 返回收到数据了
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
	addr := fmt.Sprintf(":%d", config.DefaultConfig.Server.HttpPort)
	r.Run(addr)
}
