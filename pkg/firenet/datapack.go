package firenet

import (
	"errors"
	"fire/pkg/data"
	"fire/pkg/fire"
	"fmt"
	"io"
	"k8s.io/klog/v2"
	"reflect"

	"fire/pkg/fireface"
)

// 消防协议默认数据头长度是27
var defaultBaseHeaderLen uint32 = 27

// DataPack 封包拆包类实例，暂时不需要成员
type DataPack struct {
	MaxPacketSize uint16
}

// NewDataPack 封包拆包实例初始化方法
func NewDataPack(maxPacketSize uint16) fireface.Packet {
	return &DataPack{
		MaxPacketSize: maxPacketSize,
	}
}

// GetHeadLen 获取包头长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	return defaultBaseHeaderLen
}

// Pack 封包方法(压缩数据)
func (dp *DataPack) Pack(binaryData []byte) ([]byte, error) {
	//根据传过来的数据组织返回应答数据.暂时不考虑特殊封包
	return binaryData, nil
}
func (dp *DataPack) Unpack(c fireface.IConnection) (*fire.FireMessage, error) {
	fm := &fire.FireMessage{}
	//读取客户端的Msg start
	startData := make([]byte, fire.StartSignLen)
	if _, err := io.ReadFull(c.GetTCPConnection(), startData); err != nil {
		klog.Errorf("read msg head error %v", err)
		return nil, err
	}
	if !fire.ValidationStartSign(startData) {
		return nil, fmt.Errorf("The data to be parsed has an incorrect start and end '%x' ", startData)
	}
	copy(fm.Start[:], startData[:fire.StartSignLen])
	klog.V(3).Infof("Unpack read startData %+v\n", startData)

	var controlLen int
	//读取客户端的Msg start
	controlData := make([]byte, 0)
	controlDataFirst := make([]byte, fire.FireControlFirstLen)
	controlFirstLen, err := io.ReadFull(c.GetTCPConnection(), controlDataFirst)
	if err != nil {
		klog.Errorf("read msg head error %v", err)
		return nil, err
	}
	controlLen = controlFirstLen
	control := &fire.Control{}
	if controlLen == 10 {
		klog.V(3).Infof("Unpack read Data First [%d][%x]", controlLen, controlDataFirst)
		if fire.ValidationEndSign(controlDataFirst) {
			copy(fm.End[:], controlDataFirst[8:])
			fm.ControlCommand = fire.Heartbeat
			fm.SetData(controlDataFirst[:8])
			klog.V(0).Infof("Unpack.fire-up.FireMessage Heart:[%+v]", fm)
			return fm, nil
		}
		controlTwoData := make([]byte, fire.FireControlTwoLen)
		controlTwoLen, err2 := io.ReadFull(c.GetTCPConnection(), controlTwoData)
		if err2 != nil {
			klog.Errorf("read msg head error %v", err2)
			return nil, err2
		}
		controlLen += controlTwoLen
		controlData = append(controlData, controlDataFirst...)
		controlData = append(controlData, controlTwoData...)
		klog.V(3).Infof("Unpack read Data Two [%d][%x]", controlLen, controlTwoData)
	}
	if controlLen == 25 {
		control, err = dp.UnpackControl(controlData)
		if err != nil {
			klog.Errorf("unpack error ", err)
			return nil, err
		}
		fm.Control = *control
	}

	//根据 dataLen 读取 data，放在msg.Data中
	var data []byte
	if control.GetDataLen() > 0 {
		data = make([]byte, control.GetDataLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
			klog.Errorf("read msg data error ", err)
			return nil, err
		}
		klog.V(3).Infof("Unpack read Data [%d][%x]", control.GetDataLen(), data)
	}

	var crc []byte
	crc = make([]byte, fire.CRCLen)
	if _, err := io.ReadFull(c.GetTCPConnection(), crc); err != nil {
		klog.Errorf("read msg data error ", err)
		return nil, err
	}
	fm.CRC = crc[0]
	var end []byte
	end = make([]byte, fire.EndSignLen)
	if _, err := io.ReadFull(c.GetTCPConnection(), end); err != nil {
		klog.Errorf("read msg data error ", err)
		return nil, err
	}
	copy(fm.End[:], end[:fire.StartSignLen])
	crcData := make([]byte, 0)
	crcData = append(crcData, controlData...)
	crcData = append(crcData, data...)
	// 校验数据CRC
	klog.V(3).Infof("Unpack CRC Data:Len:[%x]", crcData)
	klog.V(3).Infof("Unpack CRC[%x:%x]", fm.CRC, fire.CRC(crcData))
	if fm.CRC != fire.CRC(crcData) {
		return nil, fmt.Errorf("The CRC check fails. Procedure. Data:[%x]CRC[%x]:[%x] ", crcData, fire.CRC(crcData), fm.CRC)
	}

	fm.SetData(data)
	klog.V(3).Infof("Unpack Fire:[%+v]", fm)
	return fm, nil
}

// UnpackHeader 拆包包头方法(解压数据)
func (dp *DataPack) UnpackControl(binaryData []byte) (*fire.Control, error) {

	// 2 先处理报文前27个字节
	var control = &fire.Control{}
	typ := reflect.ValueOf(control).Type()
	var v reflect.Value
	switch typ.Kind() {
	case reflect.Ptr:
		v = reflect.New(typ.Elem()) // typ: *struct, v: *struct
	case reflect.Slice:
		v = reflect.New(typ) // typ: []x, v: *[]x
	default:
		klog.Errorf("%T is not a pointer or a slice", control)
	}
	//data.NewFireBuffer(c.Bytes).Bytes()
	if _, err := data.Decode(binaryData[:fire.FireControlLen], v.Interface()); err != nil {
		klog.Errorf("description Failed to parse the header of the fire-up packet：%v", err)
		return nil, err
	}
	klog.V(4).Infof("fire-up.Control%+v", v.Interface())

	control = v.Interface().(*fire.Control)
	if control == nil {
		return nil, fmt.Errorf("The data to be parsed is incomplete and discarded EOF")
	}
	//判断dataLen的长度是否超出我们允许的最大包长度
	if dp.MaxPacketSize > 0 && control.AppDataLen.AppDataLen > dp.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}
	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return control, nil
}
