/*
Copyright 2022 QuanxiangCloud Authors
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

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	e "fire/pkg/client/error"
	"fire/pkg/client/resp"
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"reflect"
	"time"
)

// Config client config
type Config struct {
	Timeout      time.Duration
	MaxIdleConns int
}

// New new a http client
func New(conf Config) http.Client {
	return http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(conf.Timeout * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*conf.Timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			MaxIdleConns: conf.MaxIdleConns,
		},
	}
}

// POST http post
func POST(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}) error {
	if reflect.ValueOf(entity).Kind() != reflect.Ptr {
		return errors.New("the entity type must be a pointer")
	}

	paramByte, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(paramByte)
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected state value is 200, actually %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	//klog.Infof("---------back", string(body))
	return decomposeBody(body, entity)
}
func PUT(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}) error {
	if reflect.ValueOf(entity).Kind() != reflect.Ptr {
		return errors.New("the entity type must be a pointer")
	}

	paramByte, err := json.Marshal(params)
	fmt.Println("params:", string(paramByte))
	if err != nil {
		return err
	}

	reader := bytes.NewReader(paramByte)
	var req *http.Request
	req, err = http.NewRequest("PUT", uri, reader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	var response *http.Response
	response, err = client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected state value is 200, actually %d", response.StatusCode)
	}

	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	klog.Infof(string(body))
	return decomposeBody(body, entity)
}

//func decomposeBody(body []byte) error {
//	r := new(resp.Resp)
//	if string(body) == "message delivered" {
//		r.Code = e.Success
//		r.Data = string(body)
//		return nil
//	}
//
//	if r.Code != e.Success {
//		return r.Error
//	}
//
//	return nil
//}

// GET http get
func GET(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}) error {
	if reflect.ValueOf(entity).Kind() != reflect.Ptr {
		return errors.New("the entity type must be a pointer")
	}

	paramByte, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(paramByte)
	req, err := http.NewRequest("GET", uri, reader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected state value is 200, actually %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	klog.Infof(string(body))
	return decomposeBody(body, entity)
}

func decomposeBody(body []byte, entity interface{}) error {
	r := new(resp.Resp)
	r.Data = entity

	err := json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	if r.Code != e.Success {
		return r.Error
	}

	return nil
}
