package fire

import (
	"bytes"
	"encoding/binary"
	"fire/pkg/data"
	"k8s.io/klog/v2"
	"time"
)

// MarshalRemoteControlCommandResp 处理发送数据 ，对数据确认应答0x01
func (m *FireMessage) MarshalRemoteControlCommandResp() []byte {

	var (
		send    FireMessage
		crcData = make([]byte, 0)
	)
	//
	send = *m
	control, err := data.Encode(send.Control)
	if err != nil {
		klog.Errorf("The data to be parsed is terminator incorrect, terminator '%+v' ", send.Control)
		return nil
	}
	klog.V(4).Infof("A-Fire.Control:%+v/n", send.Control)
	klog.V(4).Infof("A:Encode:%x", control)
	// 格式化control
	crcData = append(crcData, control...)

	// 格式化data
	if send.Control.GetDataLen() > 0 {
		klog.Infof("B-Fire.Data:%+v", m.Data)
		klog.Infof("B:Encode:%x", m.Data)
		crcData = append(crcData, m.Data...)
	}
	send.CRC = CRC(crcData[:len(crcData)])
	send.End = m.End
	// 校验数据CRC
	klog.V(4).Infof("C:CRC[%x:%x]/n", m.CRC, send.CRC)
	back, err := data.Encode(send)
	if err != nil {
		klog.Errorf("The data to be parsed is terminator incorrect, terminator '%+v' ", send)
		return nil
	}

	return back
}

// MarshalControlSendDataResp 处理接受四信盒子主动上报数据 ，对数据解析0x02
func (m *FireMessage) MarshalControlSendDataResp() []byte {

	var (
		send    FireMessage
		crcData = make([]byte, 0)
	)
	//
	send = *m
	control, err := data.Encode(send.Control)
	if err != nil {
		klog.Errorf("The data to be parsed is terminator incorrect, terminator '%+v' ", send.Control)
		return nil
	}
	klog.V(4).Infof("A-Fire.Control:%+v/n", send.Control)
	klog.V(4).Infof("A:Encode:%x", control)
	// 格式化control
	crcData = append(crcData, control...)

	// 格式化data
	if send.Control.GetDataLen() > 0 {
		klog.Infof("B-Fire.Data:%+v", m.Data)
		klog.Infof("B:Encode:%x", m.Data)
		crcData = append(crcData, m.Data...)
	}
	send.CRC = CRC(crcData[:len(crcData)])
	send.End = m.End
	// 校验数据CRC
	klog.V(4).Infof("C:CRC[%x:%x]/n", m.CRC, send.CRC)
	back, err := data.Encode(send)
	if err != nil {
		klog.Errorf("The data to be parsed is terminator incorrect, terminator '%+v' ", send)
		return nil
	}

	return back
}

// MarshalControlConfirmResp 处理接受四信盒子主动上报数据的确认返回，对数据确认应答0x03
func (m *FireMessage) MarshalControlConfirmResp() []byte {

	var (
		send FireMessage
		//		objectNum    uint8
		//		dataBaseType DataBaseType
		resp = make([]byte, 0)
	)
	send.Start = m.Start
	send.Control.SerialNumber.SerialNumber = m.Control.SerialNumber.SerialNumber
	send.Control.ProtocolVersion = m.Control.ProtocolVersion
	send.Control.TimeLabels = *FireToSetTime(time.Now())
	//send.FireHeader.TimeLabels = m.FireHeader.TimeLabels
	send.Control.SourceAddress.Address = m.Control.TargetAddress.Address
	send.Control.TargetAddress.Address = m.Control.SourceAddress.Address
	send.Control.AppDataLen.AppDataLen = 0
	send.Control.ControlCommand = ControlConfirm
	control, err := data.Encode(send.Control)
	if err != nil {
		klog.Errorf("The data to be parsed is terminator incorrect, terminator '%+v' ", send.Control)
		return nil
	}
	klog.V(4).Infof("A-Fire.Control:%+v", send.Control)
	resp = append(resp, control...)
	if send.GetDataLen() > 0 {
		klog.Infof("B-Fire.Data:%+v", m.Data)
		send.Data = m.Data
		resp = append(resp, m.Data...)
	}
	klog.V(4).Infof("A-Fire.CRC-Data:%x", resp)
	send.CRC = CRC(resp[:])
	send.End = m.End

	klog.V(4).Infof("FireMessage:%+v/n", send)
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写Start
	if err := binary.Write(dataBuff, binary.LittleEndian, send.Start); err != nil {
		return nil
	}
	//写 control
	if err := binary.Write(dataBuff, binary.LittleEndian, control); err != nil {
		return nil
	}
	if send.GetDataLen() > 0 {
		klog.Infof("B-Fire.Data:%+v", send.Data)
		if err := binary.Write(dataBuff, binary.LittleEndian, send.Data); err != nil {
			return nil
		}
	}
	//写crc
	if err := binary.Write(dataBuff, binary.LittleEndian, send.CRC); err != nil {
		return nil
	}
	//写End
	if err := binary.Write(dataBuff, binary.LittleEndian, send.End); err != nil {
		return nil
	}
	return dataBuff.Bytes()
}

// MarshalControlRequestResp 主动请求查询 ，查询消防设备信息0x04
func (m *FireMessage) MarshalControlRequestResp() []byte {

	var (
		send    FireMessage
		crcData = make([]byte, 0)
	)
	//
	send = *m
	control, err := data.Encode(send.Control)
	if err != nil {
		klog.Errorf("The data to be parsed is terminator incorrect, terminator '%+v' ", send.Control)
		return nil
	}
	klog.V(4).Infof("A-Fire.Control:%+v/n", send.Control)
	klog.V(4).Infof("A:Encode:%x", control)
	// 格式化control
	crcData = append(crcData, control...)

	// 格式化data
	if send.Control.GetDataLen() > 0 {
		klog.Infof("B-Fire.Data:%+v", m.Data)
		klog.Infof("B:Encode:%x", m.Data)
		crcData = append(crcData, m.Data...)
	}
	send.CRC = CRC(crcData[:len(crcData)])
	send.End = m.End
	// 校验数据CRC
	klog.V(4).Infof("C:CRC[%x:%x]/n", m.CRC, send.CRC)
	back, err := data.Encode(send)
	if err != nil {
		klog.Errorf("The data to be parsed is terminator incorrect, terminator '%+v' ", send)
		return nil
	}

	return back
}

// MarshalControlRepudiateResp 处理接受四信盒子主动上报数据和远端下发指令的确认返回 ，否认应答0x06 暂时用不上
func (m *FireMessage) MarshalControlRepudiateResp() []byte {

	var (
		send FireMessage
		//		objectNum    uint8
		//		dataBaseType DataBaseType
		resp = make([]byte, 0)
	)
	send.Start = m.Start
	send.Control.ProtocolVersion = m.Control.ProtocolVersion
	send.Control.TimeLabels = *FireToSetTime(time.Now())
	//send.FireHeader.TimeLabels = m.FireHeader.TimeLabels
	send.Control.SourceAddress.Address = m.Control.TargetAddress.Address
	send.Control.TargetAddress.Address = m.Control.SourceAddress.Address
	send.Control.AppDataLen.AppDataLen = uint16(len(m.Data))
	send.Control.ControlCommand = ControlRepudiate
	control, err := data.Encode(send.Control)
	if err != nil {
		klog.Errorf("The data to be parsed is terminator incorrect, terminator '%+v' ", send.Control)
		return nil
	}
	klog.V(4).Infof("A-Fire.Control:%+v/n", send.Control)
	resp = append(resp, control...)
	if send.Control.GetDataLen() > 0 {
		klog.Infof("B-Fire.Data:%+v", m.Data)
		resp = append(resp, m.Data...)
	}

	send.CRC = CRC(resp[:len(resp)])
	send.End = m.End

	back, err := data.Encode(send)
	if err != nil {
		klog.Errorf("The data to be parsed is terminator incorrect, terminator '%+v' ", send)
		return nil
	}
	klog.V(4).Infof("C:%+v/n", send)
	// 校验数据CRC
	klog.V(4).Infof("fireEnder[%x:%x]", send.CRC)
	return back
}

// MarshalHeartbeatResp 心跳
func (m *FireMessage) MarshalHeartbeatResp() []byte {

	// 心跳
	klog.V(4).Infof("心跳[组包开始]")
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写Start
	if err := binary.Write(dataBuff, binary.LittleEndian, m.Start); err != nil {
		return nil
	}
	//写source
	if err := binary.Write(dataBuff, binary.LittleEndian, Source); err != nil {
		return nil
	}
	//写source
	if err := binary.Write(dataBuff, binary.LittleEndian, Ethernet); err != nil {
		return nil
	}
	//写系统类型
	if err := binary.Write(dataBuff, binary.LittleEndian, General); err != nil {
		return nil
	}
	//写End
	if err := binary.Write(dataBuff, binary.LittleEndian, m.End); err != nil {
		return nil
	}
	// 心跳
	klog.V(0).Infof("心跳[组包:%x]", dataBuff.Bytes())
	return dataBuff.Bytes()
}
func GetTimestamp() []byte {
	var out [6]byte
	now := time.Now()

	out[0] = uint8(now.Second())
	out[1] = uint8(now.Minute())
	out[2] = uint8(now.Hour())
	out[3] = uint8(now.Day())
	out[4] = uint8(now.Month())
	out[5] = uint8(now.Year())
	return out[:]
}

//   ControlReserved       = 0x00 //预留
//	 RemoteControlCommand  = 0x01 //控制命令-时间同步、下发远程控制
//   收到远端下发的控制   下发控制或者拒绝
//	 ControlSendData       = 0x02 //发送数据-发送火灾自动报警系统火灾报警、运行状态等信息
//   收到设备 上送的数据信息 返回确认
//	 ControlConfirm        = 0x03 //确认-对控制命令和发送信息的确认回答
//   收到设备 上送的确认数据信息 返回确认
//	 ControlRequest        = 0x04 //请求-查询火灾自动报警系统的火灾报警、运行状态等信息
//
//	 ControlResponse       = 0x05 //应答-返回查询的信息、上报远程应答
//	 ControlRepudiate      = 0x06 //拒绝，否定-对控制命令和发送信息的否认回答
//   收到设备 上送的拒绝数据信息 返回确认

// MarshalResp 分析收到的数据对象是什么类型，格式化为消防数据包的playload
func (m *FireMessage) MarshalResp() []byte {
	// 目前只处理四信上报数据的返回应答 其他都是主动请求下发。
	//先读取客户端的数据，再组织返回数据
	klog.Infof("MarshalResp ControlCommand = %x", m.ControlCommand)

	switch m.ControlCommand {
	case RemoteControlCommand:
		return m.MarshalControlConfirmResp()
	case ControlSendData:
		return m.MarshalControlConfirmResp()
	case ControlConfirm:
		return nil
	case ControlRequest:
		return nil
	case ControlResponse:
		return nil
	case ControlRepudiate:
		return nil
	case Heartbeat:
		return m.MarshalHeartbeatResp()
	default:
		return nil
	}
}

// Marshal 根据数据对象是什么类型，格式化为消防数据包的playload
func (m *FireMessage) Marshal() []byte {
	//ControlReserved       = 0x00 //预留
	//RemoteControlCommand  = 0x01 //控制命令-时间同步、下发远程控制
	//ControlSendData       = 0x02 //发送数据-发送火灾自动报警系统火灾报警、运行状态等信息
	//ControlConfirm        = 0x03 //确认-对控制命令和发送信息的确认回答
	//ControlRequest        = 0x04 //请求-查询火灾自动报警系统的火灾报警、运行状态等信息
	//ControlResponse       = 0x05 //应答-返回查询的信息、上报远程应答
	//ControlRepudiate      = 0x06 //拒绝，否定-对控制命令和发送信息的否认回答
	switch m.ControlCommand {
	case RemoteControlCommand:
		return m.MarshalRemoteControlCommandResp()
	case ControlSendData:
		return nil
	case ControlConfirm:
		return m.MarshalControlConfirmResp()
	case ControlRequest:
		return m.MarshalControlRequestResp()
	case ControlResponse:
		return nil
	case ControlRepudiate:
		return nil
	case Heartbeat:
		return m.MarshalHeartbeatResp()
	default:
		return nil
	}
}
