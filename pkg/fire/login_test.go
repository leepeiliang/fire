package fire

import (
	"k8s.io/klog/v2"
	"testing"
	"time"
)

func TestReadRequest(t *testing.T) {
	var login = []byte{
		0x40, 0x40, 0x00, 0x00, 0x02, 0x03, 0x3a, 0x01,
		0x0f, 0x08, 0x06, 0x16, 0x58, 0x13, 0x97, 0x10,
		0xa1, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x0a, 0x00, 0x02, 0x08, 0x01, 0x00, 0x00, 0x3a,
		0x01, 0x0f, 0x08, 0x06, 0x16, 0xa8, 0x23, 0x23,
	}

	cases := CodecTestCase{
		Name: "normal",
		Struct: &FireMessage{
			Start: [2]byte{0x40, 0x40},
			Control: Control{
				//业务流水号 (2 字节)应答方按请求包的业务流水 号返回。低字节传输在前。业务流水号是一个 2 字节的正整数，
				SerialNumber: SerialNumber{SerialNumber: 0},
				//协议版本号 (2 字节)协议版本号包含主版本号(第 5 字节)和用户版本号(第 6 字节)。主版本号为 固定值 2，用户版本号由用户自行定义。
				ProtocolVersion: ProtocolVersion{
					ProtocolVersion: [2]byte{0x02, 0x03},
				},
				//数据包的第 7~12 字节，为数据包发出的时间，具体定义见 10.2.2。
				TimeLabels: TimeLabels{
					TimeLabels: [6]byte{0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16},
				},
				//数据包的第 13~18 字节，为数据包的源地址(上位机、消防控制显示装置或火 灾自动报警设备地址)。低字节传输在前。(参考建议:若发送方为平台，默认 地址为 0x00 0x00 0x00 0x00 0x00 0x00，若发送方非平台，可将源地址设置为 ID 的前六位。ID 组成为:四信 ID(1 字节)+设备类型(1 字节)+设备编号(4 字节)+预留(2 字节))
				SourceAddress: SourceAddress{
					Address: Address{
						[6]byte{0x58, 0x13, 0x97, 0x10, 0xa1, 0xff},
					},
				},
				//数据包的第 19~24 字节，为数据包的目的地址(上位机、消防控制显示装置或 火灾自动报警设备地址)。低字节传输在前。
				TargetAddress: TargetAddress{
					Address: Address{
						[6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
				},
				//数据包的第 25、26 字节，为应用数据单元的长度，长度不应大于 1024;低字节 传输在前。
				AppDataLen: AppDataLen{
					AppDataLen: 10,
				},
				//数据包的第 27 字节，控制单元的命令字节，具体定义见表 2。
				ControlCommand: 0x02,
			},

			Data: []byte{0x08, 0x01, 0x00, 0x00, 0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16},

			CRC: 0xa8,
			End: [2]byte{0x23, 0x23},
		},
		Bytes: login,
	}

	RunCodecDecodeTest(t, cases)
}

func TestReadHearderRequest(t *testing.T) {
	var login = []byte{
		0x40, 0x40, 0x00, 0x00, 0x02, 0x03, 0x3a, 0x01,
		0x0f, 0x08, 0x06, 0x16, 0x58, 0x13, 0x97, 0x10,
		0xa1, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x0a, 0x00, 0x02, 0x08, 0x01, 0x00, 0x00, 0x3a,
		0x01, 0x0f, 0x08, 0x06, 0x16, 0xa8, 0x23, 0x23,
	}

	cases := CodecTestCase{
		Name: "normal",
		Struct: &FireMessage{
			Start: [2]byte{0x40, 0x40},
			Control: Control{
				//业务流水号 (2 字节)应答方按请求包的业务流水 号返回。低字节传输在前。业务流水号是一个 2 字节的正整数，
				SerialNumber: SerialNumber{SerialNumber: 0},
				//协议版本号 (2 字节)协议版本号包含主版本号(第 5 字节)和用户版本号(第 6 字节)。主版本号为 固定值 2，用户版本号由用户自行定义。
				ProtocolVersion: ProtocolVersion{
					ProtocolVersion: [2]byte{0x02, 0x03},
				},
				//数据包的第 7~12 字节，为数据包发出的时间，具体定义见 10.2.2。
				TimeLabels: TimeLabels{
					TimeLabels: [6]byte{0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16},
				},
				//数据包的第 13~18 字节，为数据包的源地址(上位机、消防控制显示装置或火 灾自动报警设备地址)。低字节传输在前。(参考建议:若发送方为平台，默认 地址为 0x00 0x00 0x00 0x00 0x00 0x00，若发送方非平台，可将源地址设置为 ID 的前六位。ID 组成为:四信 ID(1 字节)+设备类型(1 字节)+设备编号(4 字节)+预留(2 字节))
				SourceAddress: SourceAddress{
					Address: Address{
						[6]byte{0x58, 0x13, 0x97, 0x10, 0xa1, 0xff},
					},
				},
				//数据包的第 19~24 字节，为数据包的目的地址(上位机、消防控制显示装置或 火灾自动报警设备地址)。低字节传输在前。
				TargetAddress: TargetAddress{
					Address: Address{
						[6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
				},
				//数据包的第 25、26 字节，为应用数据单元的长度，长度不应大于 1024;低字节 传输在前。
				AppDataLen: AppDataLen{
					AppDataLen: 10,
				},
				//数据包的第 27 字节，控制单元的命令字节，具体定义见表 2。
				ControlCommand: 0x02,
			},
			End: [2]byte{0x23, 0x23},
		},

		Bytes: login[:],
	}

	RunCodecHearderTest(t, cases)
}

func TestUnmarshalLogin(t *testing.T) {
	var login = []byte{
		0x40, 0x40,
		0x00, 0x00,
		0x02, 0x03,
		0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16,
		0x58, 0x13, 0x97, 0x10, 0xa1, 0xff,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x0a, 0x00,
		0x02,
		0x08, 0x01, 0x00, 0x00, 0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16,
		0xa8,
		0x23, 0x23,
	}
	//var loginback = []byte{
	//	0x40, 0x40,
	//	0x00, 0x00,
	//	0x02, 0x03,
	//	0x39, 0x22, 0x0B, 0x0C, 0x08, 0x16,
	//	0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	//	0x58, 0x13, 0x97, 0x10, 0xA1, 0xFF,
	//	0x00, 0x00,
	//	0x03,
	//}

	var fireMessage = &FireMessage{}

	klog.Infof("FireMessage:TimeLabels:%x", *FireToSetTime(time.Now()))
	klog.Infof("FireMessage:Err:%v", fireMessage.UnmarshalRes(login))
	back := fireMessage.MarshalResp()
	klog.Infof("FireMessage:Resp:%v", back)
	klog.Infof("FireMessage:%v", fireMessage)
	klog.Infof("FireMessage.TimeLabels:%v", fireMessage.Control.TimeLabels.FireToTime())
	back = fireMessage.MarshalControlSendDataResp()
	klog.Infof("FireMessage:SendDataResp%v", back)
}
