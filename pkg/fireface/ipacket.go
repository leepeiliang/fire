package fireface

import "fire/pkg/fire"

type Packet interface {
	UnpackControl(binaryData []byte) (*fire.Control, error)
	Unpack(c IConnection) (*fire.FireMessage, error)
	Pack(binaryData []byte) ([]byte, error)
	GetHeadLen() uint32
}
