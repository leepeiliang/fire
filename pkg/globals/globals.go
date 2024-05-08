/*
Copyright 2021 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package globals

import (
	"context"
	"encoding/json"
	"fire/config"
	mappercommon "fire/pkg/common"
	"fmt"
	"github.com/jasonlvhit/gocron"
	syshost "github.com/shirou/gopsutil/v3/host"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

var (
	oneHeart            *fireHeartToSouth
	MqttSubscribeClient mappercommon.MqttClient
	MqttPublishClient   mappercommon.MqttClient
)

// GetHostNameInfo 主机名称信息
func GetHostNameInfo() string {
	hostInfo, err := syshost.Info()
	if err != nil {
		klog.V(1).Infof("getHostInfo-Err:[%s]", err.Error())
		return ""
	}
	return hostInfo.Hostname
}

const (
	BJM3HostName         = "msantwo"           //1 北京M3数据中心_A
	SHWAIGAOQIAOHostName = "edge-bacnet-2"     //2 上海外高桥数据中心_A
	SANHE9               = "sanheedge9"        //3 三河9#数据中心
	TAICANG              = "taictwo"           //4 太仓数据中心_A
	NANTONGA             = "ntone"             //5 南通基地数据中心_A
	NANTONGB             = "nttwo"             //6 南通基地数据中心_B
	NANTONGC             = "ntthree"           //7 南通基地数据中心_C
	NANTONGD             = "ntfour"            //8 南通基地数据中心_D
	NANTONGE             = "ntfive"            //9 南通基地数据中心_E
	SICHUANGUANYUAN      = "sichuanguangyuana" //10 四川广元数据中心_A
	SHENZHENGARDEN       = "szgardentwo"       //11 深圳花园城数据中心_A
	JiyunHostName        = "shjiyunccc"        //12 上海纪蕴数据中心_A
	BEIJINGB28           = "btobaone"          //13 北京 B28 数据中心_A
	M6v3floors12HostName = "m6v3floors12"      //14 M5二期及三期数据中心_B
	SONGJIANG            = "shsjone"           //15 松江数据中心_A
	TongJiHostName       = "tongjithree"       //17 同济数据中心_A
	BJTONGA              = "tongzhoua"         //16 通州一号数据中心_B
	XIANJKHostName       = "xianjkccc"         //18 西安数据中心_A
	SANHE4               = "sanhe4"            //19 三河4#数据中心
	SANHE6               = "sanheedge6"        //20 三河6#数据中心
	SANHE8               = "sanheedge8"        //21 三河8#数据中心
	BJXIANGSHAN          = "bjxiangshan"       //22 北京香山数据中心
	BJM6V                = "m6v-a"             //23  M6V 数据中心_A
	BJSTALL              = "starall"           //24 星光数据中心_A
	BJSTALLb             = "starallb"          //24 星光数据中心_A
	BJSTALLc             = "starallc"          //24 星光数据中心_A
	BoXingHostName       = "boxing-bacnet"     //27 博兴数据中心_A
	SUZHOUA              = "suzhoua"           //28 宿州数据中心_A
	XIASHAHostName       = "xiasha-a"          //29 下沙数据中心_A
	NINGBOGAOXINQU       = "ningbogaoxina"     //30 宁波高新区数据中心_A
	ZHOULVSHU            = "zongshua"          //31  棕树数据中心_A
	BJM6                 = "m6a"               //32 M6数据中心_A
	BJTONGB              = "tongzhoub"         //33 通州一号数据中心_B
	TongZhouB            = "tongzhoubedge"     //33 通州一号数据中心_B
	SongJiang            = "shsjtwo"           //15 松江数据中心
	SoftwareparkHostName = "softwarepark"      //34 软件园数据中心_A
	GZFOUSHAN            = "foshan-a"          //35 佛山数据中心_A
	HEDAN                = "hedan-bacnet"      //36 荷丹数据中心_A
	YANGGUANYUN          = "yangguangyuna"     //37 央广云数据中心_A
	SHIBEI               = "shibei"            //38 市北数据中心
	JINGANG              = "jingang-a"         //39 金港数据中心_A
	TengrenHostName      = "tengren-a"         //40 腾仁数据中心_A
	SANHE7               = "sanhe7"            //20 三河7#数据中心
	WULAN4B1             = "wulan4-b1"         //51 乌兰4号数据中心_B1
	WULAN3A1             = "wulan3-a1"         //51 乌兰4号数据中心_B1
	TEST                 = "lipeiliangdeMacBook-Pro.local"
)

var defaultReSetCode = "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#"
var ReSetCode = map[string]string{
	M6v3floors12HostName: "9-#-#-#-0-1-#-#-#-#-#-#-#-#-0-#",
	//	SoftwareparkHostName: "9-#-#-#-0-0-#-#-#-#-000000-#-#-#-0-#",
	//	TengrenHostName:      "10-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	BoXingHostName:  "10-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	TengrenHostName: "9-#-#-#-0-0-#-#-#-#-0000000-#-#-#-7-#",
	SHENZHENGARDEN:  "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	JiyunHostName:   "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	TongJiHostName:  "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	XIANJKHostName:  "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	BJM3HostName:    "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	//	SHWAIGAOQIAOHostName: "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",

	TongZhouB:            "9-#-#-#-0-0-#-#-#-#-2200000-#-#-#-0-#",
	SongJiang:            "9-#-#-#-0-0-#-#-#-#-0000000-#-#-#-0-#",
	SANHE4:               "9-#-#-#-0-1-#-#-#-#-#-#-#-#-0-#",
	SANHE6:               "9-#-#-#-0-1-#-#-#-#-#-#-#-#-0-#",
	SANHE7:               "9-#-#-#-0-1-#-#-#-#-#-#-#-#-0-#",
	BEIJINGB28:           "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	GZFOUSHAN:            "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	BJM6:                 "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	TEST:                 "9-#-#-#-0-0-#-#-#-#-#-#-#-#-0-#",
	XIASHAHostName:       "1-1-0-#-#-#-#-#-#-#-#-#-#-#-13-#",
	ZHOULVSHU:            "1-1-0-#-#-#-#-#-#-#-#-#-#-#-13-#",
	SANHE8:               "1-1-0-#-#-#-#-#-#-#-#-#-#-#-13-#",
	TAICANG:              "1-1-0-#-#-#-#-#-#-#-#-#-#-#-13-#",
	SoftwareparkHostName: "1-1-0-#-#-#-#-#-#-#-#-#-#-#-13-#",
}
var HeartCode = map[string]string{
	TEST:                 "0.1",
	BJM3HostName:         "1.1",  //1 北京M3数据中心_A +
	SHWAIGAOQIAOHostName: "2.1",  //2 上海外高桥数据中心_A +
	SANHE9:               "3.1",  //3 三河9#数据中心 +
	TAICANG:              "4.1",  //4 太仓数据中心_A +
	NANTONGA:             "5.1",  //5 南通基地数据中心_A +
	NANTONGB:             "6.1",  //6 南通基地数据中心_B +
	NANTONGC:             "7.1",  //7 南通基地数据中心_C +
	NANTONGD:             "8.1",  //8 南通基地数据中心_D +
	NANTONGE:             "9.1",  //9 南通基地数据中心_E +
	SICHUANGUANYUAN:      "10.1", //10 四川广元数据中心_A +
	SHENZHENGARDEN:       "11.1", //11 深圳花园城数据中心_A +
	JiyunHostName:        "12.1", //12 上海纪蕴数据中心_A +
	BEIJINGB28:           "13.1", //13 北京 B28 数据中心_A +
	M6v3floors12HostName: "14.1", //14 M5二期及三期数据中心_B +
	SongJiang:            "15.1", //15 松江数据中心_A +
	BJTONGA:              "16.1", //33 通州一号数据中心_A
	TongJiHostName:       "17.1", //17 同济数据中心_A +
	XIANJKHostName:       "18.1", //18 西安数据中心_A, +
	SANHE4:               "19.1", //19 三河4#数据中心 +
	SANHE6:               "20.1", //20 三河6#数据中心 +
	SANHE8:               "21.1", //21 三河8#数据中心 +
	BJXIANGSHAN:          "22.1", //22 北京香山数据中心 +
	BJM6V:                "23.1", //23  M6V 数据中心_A +
	BJSTALL:              "24.1", //24 星光数据中心_A +
	BJSTALLb:             "25.1", //25 星光数据中心_b +
	BJSTALLc:             "26.1", //26 星光数据中心_c +
	BoXingHostName:       "27.1", //27 博兴数据中心_A +
	SUZHOUA:              "28.1", //28 宿州数据中心_A +
	XIASHAHostName:       "29.1", //29 下沙数据中心_A +
	NINGBOGAOXINQU:       "30.1", //30 宁波高新区数据中心_A +
	ZHOULVSHU:            "31.1", //31  棕树数据中心_A +
	BJM6:                 "32.1", //32 M6数据中心_A +
	BJTONGB:              "33.1", //33 通州一号数据中心_B +
	TongZhouB:            "33.1", //33 通州一号数据中心_B +
	SoftwareparkHostName: "34.1", //34 软件园数据中心_A +
	GZFOUSHAN:            "35.1", //35 佛山数据中心_A +
	HEDAN:                "36.1", //36 荷丹数据中心_A +
	YANGGUANYUN:          "37.1", //37 央广云数据中心_A +
	SHIBEI:               "38.1", //38 市北数据中心_A +
	JINGANG:              "39.1", //39 金港数据中心_A
	TengrenHostName:      "40.1", //40 腾仁数据中心_A +
	SANHE7:               "41.1", //41 三河7#数据中心

	WULAN4B1: "51.1", //51 乌兰4号数据中心_B1
	WULAN3A1: "43.1", //51 乌兰4号数据中心_B1
}
var HeartModelID = "1.4.1.1"       // 2.3.1.1
var HeartPropertyID = "1.5.9998.1" // 1.2.1.1

const DOT = "."

// GetHostReSetCodeInfo 主机名称信息
func GetHostReSetCodeInfo() string {
	hostInfo, err := syshost.Info()
	if err != nil {
		klog.V(1).Infof("getHostInfo-Err:[%s]", err.Error())
		return ""
	}
	klog.Infof("HostName:[%s]", hostInfo.Hostname)
	if code, ok := ReSetCode[hostInfo.Hostname]; ok {
		//		klog.Infof("ReSetCode:[%s]", code)
		return code
	}
	return defaultReSetCode
}

// GetHostHeartDeviceID 各机房心跳点位
func GetHostHeartDeviceID() string {
	hostInfo, err := syshost.Info()
	if err != nil {
		klog.V(1).Infof("getHostInfo-Err:[%s]", err.Error())
		return ""
	}
	klog.Infof("DeviceID-HostName:[%s]", hostInfo.Hostname)
	if code, ok := HeartCode[hostInfo.Hostname]; ok {

		property := []string{code, HeartModelID}
		klog.Infof("DeviceID:[%s]", property)
		return strings.Join(property, DOT)
	}
	return "0.0"
}

// GetHostHeartPropertyID 各机房心跳点位
func GetHostHeartPropertyID() string {
	hostInfo, err := syshost.Info()
	if err != nil {
		klog.V(1).Infof("getHostInfo-Err:[%s]", err.Error())
		return ""
	}
	klog.Infof("PropertyID-HostName:[%s]", hostInfo.Hostname)
	if code, ok := HeartCode[hostInfo.Hostname]; ok {

		property := []string{code, HeartModelID, HeartPropertyID}
		klog.Infof("PropertyID:[%s]", property)
		return strings.Join(property, DOT)
	}
	return "0.0"
}

// GetHostTopicInfo 推送消息到南向的topic信息
func GetHostTopicInfo() string {
	hostInfo, err := syshost.Info()
	if err != nil {
		klog.V(1).Infof("getHostInfo-Err:[%s]", err.Error())
		return ""
	}
	klog.Infof("HostName:[%s]", hostInfo.Hostname)

	return fmt.Sprintf(config.DefaultConfig.Topic.DeviceUpdateData, hostInfo.Hostname)
}

// FireHeartToSouth FireHeartToSouth
type FireHeartToSouth interface {
	HeartProperties(in int)
	ReSetSenHeart() error
}

// fireHeartToSouth 同步更新角色和用户关系数据
type fireHeartToSouth struct {
	Ctx          context.Context
	Heart        *gocron.Scheduler
	SendTimeInfo chan *TimeInfo
}

// TimeInfo 时间戳
type TimeInfo struct {
	SigTime time.Time
}

// New 初始化一个心跳管道和定时器
func New() {
	oneHeart = &fireHeartToSouth{
		Ctx:   context.Background(),
		Heart: gocron.NewScheduler(),
	}
	oneHeart.ReSetSenHeart()
}

// GetfireHeart  获取心跳
func GetfireHeart() FireHeartToSouth {
	if oneHeart == nil {
		return nil
	}
	return oneHeart
}

func (s *fireHeartToSouth) HeartProperties(in int) {

	var (
		err     error
		topic   = GetHostTopicInfo()
		dataMsg mappercommon.DeviceGroupCustomizedData
	)
	klog.V(4).Infof("心跳点位组包")

	var tmp = make([]mappercommon.DeviceCustomizedData, 0)
	var updateMsg mappercommon.DeviceCustomizedData
	updateMsg.Timestamp = mappercommon.GetTimestamp()
	updateMsg.Data = map[string]*mappercommon.DataValue{}
	dataMsg.DeviceType = HeartModelID
	updateMsg.DeviceID = GetHostHeartDeviceID()
	dataMsg.Timestamp = mappercommon.GetTimestamp()
	updateMsg.Data[GetHostHeartPropertyID()] = &mappercommon.DataValue{
		Value:     in,
		Timestamp: mappercommon.GetTimestamp(),
		Metadata: mappercommon.DataMetadata{
			Type:      "boolean",
			Timestamp: mappercommon.GetTimestamp(),
		},
	}
	tmp = append(tmp, updateMsg)
	dataMsg.Devices = tmp
	var payload []byte

	if payload, err = json.Marshal(dataMsg); err != nil {
		klog.Error("Create message data failed")
		return
	}
	klog.V(3).Infof("Device Property[%s]", string(payload))
	if err = MqttPublishClient.Publish(topic, payload); err != nil {
		klog.Error(err)
		return
	}

	klog.V(2).Infof("Update value: %s, topic: %s", dataMsg.DeviceType, topic)

}
func (s *fireHeartToSouth) ReSetSenHeart() error {
	now := time.Now()
	s.Heart.Clear()
	s.Heart.ChangeLoc(time.UTC)
	job := s.Heart.Every(920).Second()
	job.Do(s.HeartProperties, 1)
	next := job.NextScheduledTime()
	expected := now.UTC().Add(920 * time.Second)
	fmt.Println(expected)
	fmt.Println(next.UTC())
	s.Heart.Start()
	return nil

}
