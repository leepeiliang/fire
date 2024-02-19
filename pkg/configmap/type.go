package configmap

// Common visitor configurations for BacNet protocol
type VisitorConfig struct {
	// Required: name of customized protocol
	ProtocolName string `json:"protocolName,omitempty"`
	// Required: The configData of customized protocol
	ConfigData CustomizedValue `json:"configData,omitempty"`
}

// ProtocolCommonConfig is the BacNet protocol configuration.
type ProtocolCommonConfig struct {
	CustomizedValues CustomizedValue `json:"customizedValues,omitempty"`
}

// ProtocolConfig is the protocol configuration.
type ProtocolConfig struct {
	// Unique protocol name
	// Required.
	ProtocolName string `json:"protocolName,omitempty"`
	// Any mqttconfig data
	// +optional
	ConfigData CustomizedValue `json:"configData,omitempty"`
}

// CustomizedValue is the customized part for bacnet protocol.
type CustomizedValue map[string]interface{}
