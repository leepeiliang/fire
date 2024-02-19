package firenet

import (
	"fire/pkg/fire"
	"fmt"
	"io"
	"k8s.io/klog/v2"
	"net"
	"testing"
	"time"
)

// run in terminal:
// go test -v ./firenet -run=TestDataPack

// 只是负责测试datapack拆包，封包功能
func TestDataPack(t *testing.T) {
	//创建socket TCP Server
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建服务器gotoutine，负责从客户端goroutine读取粘包的数据，然后进行解析
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				klog.Errorf("server accept err:", err)
			}

			//处理客户端请求
			go func(conn net.Conn) {
				//创建封包拆包对象dp
				dp := NewDataPack(2048)
				for {
					//msg 是有data数据的，需要再次读取data数据
					msg := fire.FireMessage{}
					//1 先读出流中的head部分
					startData := make([]byte, fire.StartSignLen)
					_, err := io.ReadFull(conn, startData) //ReadFull 会把msg填充满为止
					if err != nil {
						klog.Errorf("read head error")
						return
					}
					copy(msg.Start[:], startData[:fire.StartSignLen])
					//2 先读出流中的控制部分
					controData := make([]byte, fire.FireControlLen)
					_, err = io.ReadFull(conn, controData) //ReadFull 会把msg填充满为止
					if err != nil {
						klog.Errorf("read head error")
						return
					}
					//将Control字节流 拆包到msg中
					msgControl, err := dp.UnpackControl(controData)
					if err != nil {
						klog.Errorf("server unpack err:", err)
						return
					}
					msg.Control = *msgControl
					klog.Infof("==> Recv msg:start=%x, Control=%x", startData, controData)
					klog.Infof("==> Recv msg:struct=%v", msg.Control)
					msg.Data = make([]byte, msgControl.GetDataLen())
					if msg.Control.GetDataLen() > 0 {
						//根据dataLen从io中读取字节流
						backlen, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							klog.Errorf("server unpack data err:", err)
							return
						}

						klog.Infof("==> Recv DataMsg:len=%d, data=%x backlen=%d", msgControl.GetDataLen(), msg.Data, backlen)
					}
					var crc []byte
					crc = make([]byte, fire.CRCLen)
					_, err = io.ReadFull(conn, crc)
					if err != nil {
						klog.Errorf("read msg data error ", err)
						return
					}
					msg.CRC = crc[0]
					klog.Infof("==> Recv DataMsg:crc=%x", crc)
					var end []byte
					end = make([]byte, fire.EndSignLen)
					if _, err := io.ReadFull(conn, end); err != nil {
						klog.Errorf("read msg data error ", err)
						return
					}
					copy(msg.End[:], end[:fire.EndSignLen])

					klog.Infof("==> Recv EndData:crc=%x, data=%x", msg.CRC, msg.End)
				}
			}(conn)
		}
	}()

	//客户端goroutine，负责模拟粘包的数据，然后进行发送
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:7777")
		if err != nil {
			klog.Errorf("client dial err:", err)
			return
		}

		//创建一个封包对象 dp
		dp := NewDataPack(2048)

		//封装一个msg1包
		msg1 := &fire.FireMessage{

			Start: [2]byte{0x40, 0x40},
			Control: fire.Control{
				//业务流水号 (2 字节)应答方按请求包的业务流水 号返回。低字节传输在前。业务流水号是一个 2 字节的正整数，
				SerialNumber: fire.SerialNumber{SerialNumber: 0},
				//协议版本号 (2 字节)协议版本号包含主版本号(第 5 字节)和用户版本号(第 6 字节)。主版本号为 固定值 2，用户版本号由用户自行定义。
				ProtocolVersion: fire.ProtocolVersion{
					ProtocolVersion: [2]byte{0x02, 0x03},
				},
				//数据包的第 7~12 字节，为数据包发出的时间，具体定义见 10.2.2。
				TimeLabels: fire.TimeLabels{
					TimeLabels: [6]byte{0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16},
				},
				//数据包的第 13~18 字节，为数据包的源地址(上位机、消防控制显示装置或火 灾自动报警设备地址)。低字节传输在前。(参考建议:若发送方为平台，默认 地址为 0x00 0x00 0x00 0x00 0x00 0x00，若发送方非平台，可将源地址设置为 ID 的前六位。ID 组成为:四信 ID(1 字节)+设备类型(1 字节)+设备编号(4 字节)+预留(2 字节))
				SourceAddress: fire.SourceAddress{
					Address: fire.Address{
						[6]byte{0x58, 0x13, 0x97, 0x10, 0xa1, 0xff},
					},
				},
				//数据包的第 19~24 字节，为数据包的目的地址(上位机、消防控制显示装置或 火灾自动报警设备地址)。低字节传输在前。
				TargetAddress: fire.TargetAddress{
					Address: fire.Address{
						[6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
				},
				//数据包的第 25、26 字节，为应用数据单元的长度，长度不应大于 1024;低字节 传输在前。
				AppDataLen: fire.AppDataLen{
					AppDataLen: 10,
				},
				//数据包的第 27 字节，控制单元的命令字节，具体定义见表 2。
				ControlCommand: 0x02,
			},

			Data: []byte{0x08, 0x01, 0x00, 0x00, 0x3a, 0x01, 0x0f, 0x08, 0x06, 0x16},

			CRC: 0xa8,
			End: [2]byte{0x23, 0x23},
		}
		sendData1, err := dp.Pack(msg1.MarshalControlSendDataResp())
		if err != nil {
			klog.Errorf("client pack msg1 err:", err)
			return
		}
		klog.Infof("==> Send data=%x", sendData1)
		//msg2 := &fire-up.FireMessage{}
		//sendData2, err := dp.Pack(msg2)
		//if err != nil {
		//	fmt.Println("client temp msg2 err:", err)
		//	return
		//}
		//
		////将sendData1，和 sendData2 拼接一起，组成粘包
		//sendData1 = append(sendData1, sendData2...)

		//向服务器端写数据
		conn.Write(sendData1)
	}()

	//客户端阻塞
	select {
	case <-time.After(time.Second):
		return
	}
}
