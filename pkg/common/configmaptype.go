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

package common

import (
	"encoding/json"
	"time"
)

// DeviceProfile is structure to store in configMap.
type DeviceProfile struct {
	DeviceInstances []DeviceInstance `json:"deviceInstances,omitempty"`
	DeviceModels    []DeviceModel    `json:"deviceModels,omitempty"`
	Protocols       []Protocol       `json:"protocols,omitempty"`
}

// DeviceInstance is structure to store device in deviceProfile.json.bak in configmap.
type DeviceInstance struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	ProtocolName string `json:"protocol,omitempty"`
	PProtocol    Protocol
	Model        string `json:"model,omitempty"`
	Twins        []Twin `json:"twins,omitempty"`
	//	Properties   []DataProperty `json:"dataProperties,omitempty"`
	//	Topic        string         `json:"datatopic,omitempty"`
	DData            Data              `json:"data,omitempty"`
	PropertyVisitors []PropertyVisitor `json:"propertyVisitors,omitempty"`
}

// DeviceModel is structure to store deviceModel in deviceProfile.json.bak in configmap.
type DeviceModel struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Properties  []Property `json:"properties,omitempty"`
}

// Property is structure to store deviceModel property.
type Property struct {
	Name         string      `json:"name,omitempty"`
	DataType     string      `json:"dataType,omitempty"`
	Description  string      `json:"description,omitempty"`
	AccessMode   string      `json:"accessMode,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Minimum      interface{} `json:"minimum,omitempty"`
	Maximum      interface{} `json:"maximum,omitempty"`
	Unit         string      `json:"unit,omitempty"`
}

// Protocol is structure to store protocol in deviceProfile.json.bak in configmap.
type Protocol struct {
	Name                 string          `json:"name,omitempty"`
	Protocol             string          `json:"protocol,omitempty"`
	ProtocolConfigs      json.RawMessage `json:"protocolConfig,omitempty"`
	ProtocolCommonConfig json.RawMessage `json:"protocolCommonConfig,omitempty"`
}

// PropertyVisitor is structure to store propertyVisitor in deviceProfile.json.bak in configmap.
type PropertyVisitor struct {
	Name             string `json:"name,omitempty"`
	PropertyName     string `json:"propertyName,omitempty"`
	ModelName        string `json:"modelName,omitempty"`
	CollectCycle     int64  `json:"collectCycle"`
	ReportCycle      int64  `json:"reportcycle,omitempty"`
	PProperty        Property
	Protocol         string          `json:"protocol,omitempty"`
	CustomizedValues CustomizedValue `json:"customizedValues,omitempty"`
	VisitorConfig    json.RawMessage `json:"visitorConfig"`
}

type CustomizedValue map[string]interface{}

// // Data is data structure for the message that only be subscribed in edge node internal.
type Data struct {
	Properties []DataProperty `json:"dataProperties,omitempty"`
	Topic      string         `json:"datatopic,omitempty"`
}

// DataProperty is data property.
type DataProperty struct {
	MetaData     map[string]string `json:"metadata,omitempty"`
	PropertyName string            `json:"propertyName,omitempty"`
	PVisitor     *PropertyVisitor
}

// Metadata is the metadata for data.
type Metadata struct {
	Timestamp string `json:"timestamp,omitempty"`
	Type      string `json:"type,omitempty"`
}

// Twin is the set/get pair to one register.
type Twin struct {
	PropertyName string `json:"propertyName,omitempty"`
	PVisitor     *PropertyVisitor
	Desired      DesiredData  `json:"desired,omitempty"`
	Reported     ReportedData `json:"reported,omitempty"`
}

// DesiredData is the desired data.
type DesiredData struct {
	Value     string   `json:"value,omitempty"`
	Metadatas Metadata `json:"metadata,omitempty"`
}

// ReportedData is the reported data.
type ReportedData struct {
	Value     string   `json:"value,omitempty"`
	Metadatas Metadata `json:"metadata,omitempty"`
}

type BaseDevice struct {
	ID             int64         `json:"id"                       gorm:"primaryKey;autoIncrement"`
	HostName       string        `json:"hostName,omitempty"       gorm:"column:host_name"`
	NodeName       string        `json:"nodeName,omitempty"       gorm:"column:node_name"`
	Protocol       string        `json:"protocol,omitempty"       gorm:"column:protocol"`
	DeviceID       string        `json:"deviceID,omitempty"       gorm:"column:device_id"`
	DeviceName     string        `json:"deviceName,omitempty"     gorm:"column:device_name"`
	ProtocolConfig string        `json:"protocolConfig,omitempty" gorm:"column:protocolConfig"`
	BrandModel     string        `json:"brandModel,omitempty"     gorm:"column:brand_model"`
	Status         string        `json:"status,omitempty"         gorm:"column:status"`
	ModelType      string        `json:"modelType,omitempty"      gorm:"column:model_type"`
	Location       string        `json:"location,omitempty"       gorm:"column:location"`
	CreatedAt      time.Time     `json:"createdAt,omitempty"      gorm:"column:created_at"`
	UpdatedAt      time.Time     `json:"updatedAt,omitempty"      gorm:"column:updated_at"`
	CreatedUser    string        `json:"createdUser,omitempty"    gorm:"column:created_user"`
	UpdatedUser    string        `json:"updatedUser,omitempty"    gorm:"column:updated_user"`
	Properties     []PropertyNew `json:"properties,omitempty"     gorm:"foreignKey:AssociationID"`
}

type PropertyNew struct {
	AssociationID int64  `json:"AssociationID,omitempty" gorm:"column:association_id"`
	PropertyName  string `json:"propertyName,omitempty"  gorm:"column:property_name"`
	PropertyID    string `json:"propertyID,omitempty"    gorm:"column:property_id"`
	DataType      string `json:"dataType,omitempty"      gorm:"column:data_type"`
	AccessMode    string `json:"accessMode,omitempty"    gorm:"column:access_mode"`
	VisitorConfig string `json:"visitorConfig,omitempty" gorm:"column:visitor_config"`
	Active        string `json:"active,omitempty"        gorm:"column:active"`
	Unit          string `json:"unit,omitempty"          gorm:"column:unit"`
}
