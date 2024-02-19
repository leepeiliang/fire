// Copyright 2018-2020 opcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package data

import (
	"fmt"
	"k8s.io/klog/v2"
	"math"
	"reflect"
	"time"
)

var (
	binaryDecoder = reflect.TypeOf((*BinaryDecoder)(nil)).Elem()
	timeType      = reflect.TypeOf(time.Time{})
)

func isBinaryDecoder(val reflect.Value) bool {
	return val.Type().Implements(binaryDecoder)
}

func isTime(val reflect.Value) bool {
	return val.Type() == timeType
}

type BinaryDecoder interface {
	Decode([]byte) (int, error)
}

func Decode(b []byte, v interface{}) (int, error) {
	val := reflect.ValueOf(v)
	return decode(b, val, val.Type().String())
}

func decode(b []byte, val reflect.Value, name string) (n int, err error) {
	if debugCodec {
		fmt.Printf("decode: %s has type %v and is a %s, %d bytes\n", name, val.Type(), val.Type().Kind(), len(b))
		defer func() {
			fmt.Printf("decode: decoded %d bytes into %s\n", n, name)
		}()
	}
	defer klog.V(4).Infof("decode: %s:%v\n", name, val)
	buf := NewBuffer(b)

	klog.V(4).Infof("buf:%x/n", buf)
	switch {
	case isBinaryDecoder(val):
		v := val.Interface().(BinaryDecoder)
		return v.Decode(b)
	case isTime(val):
		val.Set(reflect.ValueOf(buf.ReadTime()))
	default:
		klog.V(4).Infof("decode: %s is a %s\n", name, val.Kind())
		switch val.Kind() {
		case reflect.Bool:
			val.SetBool(buf.ReadBool())
		case reflect.Int8:
			val.SetInt(int64(buf.ReadInt8()))
		case reflect.Uint8:
			val.SetUint(uint64(buf.ReadByte()))
		case reflect.Int16:
			val.SetInt(int64(buf.ReadInt16()))
		case reflect.Uint16:
			val.SetUint(uint64(buf.ReadUint16()))
		case reflect.Int32:
			val.SetInt(int64(buf.ReadInt32()))
		case reflect.Uint32:
			val.SetUint(uint64(buf.ReadUint32()))
		case reflect.Int64:
			val.SetInt(buf.ReadInt64())
		case reflect.Uint64:
			val.SetUint(buf.ReadUint64())
		case reflect.Float32:
			val.SetFloat(float64(buf.ReadFloat32()))
		case reflect.Float64:
			val.SetFloat(buf.ReadFloat64())
		case reflect.String:
			val.SetString(buf.ReadString(readUint16))
		case reflect.Slice:
			return decodeSlice(b, val, name)
		case reflect.Array:
			return decodeArray(b, val, name)
		case reflect.Ptr:
			return decode(b, val.Elem(), name)
		case reflect.Struct:
			return decodeStruct(b, val, name)
		default:
			return 0, fmt.Errorf("unsupported type %s", val.Type())
		}

	}
	return buf.Pos(), buf.Error()
}

func decodeStruct(b []byte, val reflect.Value, name string) (int, error) {
	pos := 0
	valt := val.Type()
	//	klog.Infof("val.Type:%+v", valt)
	for i := 0; i < val.NumField(); i++ {
		ft := valt.Field(i)
		fname := name + "." + ft.Name
		// if the field is a pointer we need to create
		// the value before we can marshal data into it.
		f := val.Field(i)
		//klog.Infof("valt.Field(%d):%+v----%+v", i, f, f.Type().Kind())
		if f.Type().Kind() == reflect.Ptr {
			f.Set(reflect.New(f.Type().Elem()))
			//klog.Infof("Elem----:%+v", reflect.New(f.Type().Elem()))
			klog.V(0).Infof("decode: %s has type %v and has new value %#v\n", fname, f.Type(), f.Interface())
		}
		klog.V(4).Infof("offset%d", pos)
		n, err := decode(b[pos:], f, fname)
		if err != nil {
			return pos, err
		}
		pos += n
	}
	return pos, nil
}

func decodeSlice(b []byte, val reflect.Value, name string) (int, error) {
	buf := NewBuffer(b)
	n := buf.ReadInt8()
	if buf.Error() != nil {
		klog.Errorf("%v", buf.Error())
		return buf.Pos(), buf.Error()
	}

	if uint8(n) >= null8 {
		return buf.Pos(), nil
	}

	if int(n) > math.MaxUint8 {
		return buf.Pos(), fmt.Errorf("array too large: %d > %d", n, math.MaxUint8)
	}

	// elemType is the type of the slice elements
	// e.g. *Foo for []*Foo
	elemType := val.Type().Elem()
	// fmt.Println("elemType: ", elemType.String())

	// fast path for []byte
	if elemType.Kind() == reflect.Uint8 {
		// fmt.Println("decode: []byte fast path")
		val.SetBytes(buf.ReadN(int(n)))
		return buf.Pos(), buf.Error()
	}

	pos := buf.Pos()
	// a is a slice of []*Foo
	a := reflect.MakeSlice(val.Type(), int(n), int(n))
	for i := 0; i < int(n); i++ {

		// if the slice elements are pointers we need to create
		// them before we can marshal data into them.
		if elemType.Kind() == reflect.Ptr {
			a.Index(i).Set(reflect.New(elemType.Elem()))
		}

		ename := fmt.Sprintf("%s[%d]", name, i)
		m, err := decode(b[pos:], a.Index(i), ename)
		if err != nil {
			return pos, err
		}
		pos += m
	}
	val.Set(a)

	return pos, nil
}

func decodeArray(b []byte, val reflect.Value, name string) (int, error) {
	buf := NewBuffer(b)
	//	fmt.Printf("message:Name:[%s]Type[%+v]Kind:[%+v]Valuelen:[%d]\n", name, val.Type(), val.Kind(), val.Len())

	//n := buf.ReadUint32()
	n := uint32(val.Len())
	if buf.Error() != nil {
		return buf.Pos(), buf.Error()
	}

	if uint8(n) == null8 {
		return buf.Pos(), nil
	}

	if n > math.MaxInt32 {
		return buf.Pos(), fmt.Errorf("array too large: %d > %d", n, math.MaxInt32)
	}

	if n > uint32(val.Len()) {
		return buf.Pos(), fmt.Errorf("array too large: %d > %d", n, val.Len())
	}

	// elemType is the type of the slice elements
	// e.g. *Foo for []*Foo
	elemType := val.Type().Elem()

	klog.V(4).Infof("elemType: %v", elemType.String())
	// fast path for []byte
	if elemType.Kind() == reflect.Uint8 {
		//		fmt.Printf("decode: [%x]byte fast path\n",buf.ReadN(int(n)))
		reflect.Copy(val, reflect.ValueOf(buf.ReadN(int(n))))
		return buf.Pos(), buf.Error()
	}

	pos := buf.Pos()
	// a is a pointer to an array [n]*Foo, where n is know at compile time
	a := reflect.New(val.Type()).Elem()
	for i := 0; i < int(n); i++ {

		// if the slice elements are pointers we need to create
		// them before we can marshal data into them.
		if elemType.Kind() == reflect.Ptr {
			a.Index(i).Set(reflect.New(elemType.Elem()))
		}

		ename := fmt.Sprintf("%s[%d]", name, i)
		m, err := decode(b[pos:], a.Index(i), ename)
		if err != nil {
			return pos, err
		}
		pos += m
	}
	val.Set(a)

	return pos, nil
}
