package fire

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"fire/pkg/data"

	"k8s.io/klog/v2"
)

func (m *FireMessage) UnmarshalRes(dAtA []byte) error {
	l := len(dAtA)

	iNdEx := 0
	if l < StartSignLen+EndSignLen {
		return fmt.Errorf("The data to be parsed is incomplete and discarded EOF")
	}
	if !ValidationStartSign(dAtA) {
		return fmt.Errorf("The data to be parsed has an incorrect start and end '%s' ", string(dAtA[:2]))
	}
	if !ValidationEndSign(dAtA) {
		return fmt.Errorf("The data to be parsed is terminator incorrect, terminator '%s' ", string(dAtA[l-2:l]))
	}
	// 1 处理报文尾部加盐部分
	iNdEx += 2
	copy(m.Start[:], dAtA[:iNdEx])
	klog.V(4).Infof("fireStart:%+v", m.Start)

	end := l - EndSignLen
	copy(m.End[:], dAtA[end:l])
	klog.V(4).Infof("fireEnder.End:len:[%d:%d]data:[%x]", end, l, m.End)
	end = end - CRCLen
	m.CRC = dAtA[end]
	klog.V(4).Infof("fireEnder.CRC:len:[%d:%d]data:[%x:%x]", end, l, m.CRC, CRC(dAtA[2:l-3]))
	// 校验数据CRC
	if m.CRC != CRC(dAtA[iNdEx:l-3]) {
		return fmt.Errorf("The CRC check fails. Procedure. Data:[%x]CRC[%x]:[%x] ", dAtA[2:l-3], CRC(dAtA[2:l-3]), m.CRC)
	}

	klog.V(4).Infof("fireEnder:%+v", m.End)

	// 2 先处理报文前2-27个字节,控制单元
	var control = &Control{}
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
	if _, err := data.Decode(dAtA[iNdEx:StartSignLen+FireControlLen], v.Interface()); err != nil {
		klog.Errorf("description Failed to parse the header of the fire-up packet：%v", err)
		return err
	}
	iNdEx += 25
	klog.V(4).Infof("Interface%+v", v.Interface())

	m.Control = *v.Interface().(*Control)
	if &m.Control == nil {
		return fmt.Errorf("The data to be parsed is incomplete and discarded EOF")
	}
	klog.V(4).Infof("Control%+v", m.Control)
	// 3 处理报文数据部分
	var data = make([]byte, 0)

	klog.Infof("DateLen:%d ", m.Control.AppDataLen)

	data = append(data, dAtA[iNdEx:iNdEx+int(m.Control.AppDataLen.AppDataLen)]...)
	iNdEx += int(m.Control.AppDataLen.AppDataLen)

	klog.V(4).Infof("Date:%x ", data)

	m.Data = data
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func ValidationStartSign(in []byte) bool {
	return bytes.HasPrefix(in, StartSign[:])
}

func ValidationEndSign(in []byte) bool {
	return bytes.HasSuffix(in, EndSign[:])
}
