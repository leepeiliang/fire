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

package config

import (
	"errors"
	"io/ioutil"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var DefaultConfig Config

// Config is the modbus mapper configuration.
type Config struct {
	Server     Server     `yaml:"server"`
	EdgeServer EdgeServer `yaml:"edgeserver"`
	Mqtt       Mqtt       `yaml:"mqtt"`
	Configmap  string     `yaml:"configmap"`
	DomanCode  DomanCode  `yaml:"domancode"`
}
type Server struct {
	IPVersion string `yaml:"iPVersion,omitempty"` //当前服务器协议
	Host      string `yaml:"host,omitempty"`      //当前服务器主机IP
	TCPPort   int    `yaml:"tcpPort,omitempty"`   //当前服务器主机监听端口号
	HttpPort  int    `yaml:"httpPort,omitempty"`  //当前服务器参数下发监听端口号
	Name      string `yaml:"name,omitempty"`      //当前服务器名称

	/*
		Zinx
	*/
	Version          string `yaml:"version,omitempty"`          //当前Zinx版本号
	MaxPacketSize    uint16 `yaml:"maxPacketSize,omitempty"`    //都需数据包的最大值
	MaxConn          int    `yaml:"maxConn,omitempty"`          //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint16 `yaml:"workerPoolSize,omitempty"`   //业务工作Worker池的数量
	MaxWorkerTaskLen uint16 `yaml:"maxWorkerTaskLen,omitempty"` //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint16 `yaml:"maxMsgChanLen,omitempty"`    //SendBuffMsg发送消息的缓冲最大长度
}

type EdgeServer struct {
	Host string `yaml:"host,omitempty"` //当前服务器主机IP
	Port int    `yaml:"port,omitempty"` //当前服务器参数下发监听端口号
	Name string `yaml:"name,omitempty"` //当前服务器名称
}

// Mqtt is the Mqtt configuration.
type Mqtt struct {
	ServerAddress  string `yaml:"server,omitempty"`
	UserName       string `yaml:"username,omitempty"`
	Password       string `yaml:"password,omitempty"`
	CertFile       string `yaml:"certification,omitempty"`
	PrivateKeyFile string `yaml:"privatekey,omitempty"`
	Qos            int    `yaml:"qos,omitempty"`
	Retained       bool   `yaml:"retained,omitempty"`
}

type DomanCode struct {
	Code map[string]string `yaml:"code"`
}

// ErrConfigCert error of certification configuration.
var ErrConfigCert = errors.New("Both certification and private key must be provided")

var defaultConfigFile = "/Users/lipeiliang/go/src/github.com/leepeiliang/fire/config/config.yaml"

//var defaultConfigFile = "kubeedge/etc/config.yaml"

// Parse parse the configuration file. If failed, return error.
func (c *Config) Parse() error {
	var level klog.Level
	var loglevel string
	var configFile string

	pflag.StringVar(&loglevel, "v", "3", "log level")
	pflag.StringVar(&configFile, "config-file", defaultConfigFile, "Config file name")
	pflag.Parse()
	cf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(cf, c); err != nil {
		return err
	}
	if err = level.Set(loglevel); err != nil {
		return err
	}

	return c.parseFlags()
}

// parseFlags parse flags. Certification and Private key must be provided at the same time.
func (c *Config) parseFlags() error {
	pflag.StringVar(&c.Mqtt.ServerAddress, "mqtt-address", c.Mqtt.ServerAddress, "MQTT broker address")
	pflag.StringVar(&c.Mqtt.UserName, "mqtt-username", c.Mqtt.UserName, "username")
	pflag.StringVar(&c.Mqtt.Password, "mqtt-password", c.Mqtt.Password, "password")
	pflag.StringVar(&c.Mqtt.CertFile, "mqtt-certification", c.Mqtt.CertFile, "certification file path")
	pflag.StringVar(&c.Mqtt.PrivateKeyFile, "mqtt-priviatekey", c.Mqtt.PrivateKeyFile, "private key file path")

	pflag.Parse()

	if (c.Mqtt.CertFile != "" && c.Mqtt.PrivateKeyFile == "") ||
		(c.Mqtt.CertFile == "" && c.Mqtt.PrivateKeyFile != "") {
		return ErrConfigCert
	}

	return nil
}
