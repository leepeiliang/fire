/*
Copyright 2020 The KubeEdge Authors.

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

package configmap

import (
	"bytes"
	"encoding/json"
	mappercommon "fire/pkg/common"

	"io/ioutil"

	"k8s.io/klog/v2"
)

// NewParse parse the configmap.
func NewParse(path string, devices map[string]*mappercommon.BaseDevice) error {

	klog.Infof("configmap path:%s", path)
	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		klog.Errorf("configmap Parse:%s ", err.Error())
		return err
	}
	if len(jsonFile) <= 0 {
		return nil
	}
	klog.Info("body-len:%d", len(jsonFile))
	jsonFile = bytes.TrimPrefix(jsonFile, []byte("\xef\xbb\xbf"))
	if err = json.Unmarshal(jsonFile, &devices); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	klog.Infof("deviceProfile.DeviceInstances:len:[%d]", len(devices))
	//for k,v:=range devices{
	//	klog.Infof("deviceProfile.DeviceInstances:len:[%d]", k)
	//	klog.Infof("deviceProfile.DeviceInstances:len:[%v]", v)
	//}

	return nil
}
