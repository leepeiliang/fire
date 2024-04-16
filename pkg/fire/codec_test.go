// Copyright 2018-2020 opcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// This file is copied to packages ua, uacp and uasc to break an import cycle.

package fire

import (
	"fire/pkg/data"
	"reflect"
	"testing"

	"github.com/pascaldekloe/goe/verify"
	"k8s.io/klog/v2"
)

// CodecTestCase describes a test.md case for a encoding and decoding an
// object from bytes.
type CodecTestCase struct {
	Name   string
	Struct interface{}
	Bytes  []byte
}

// RunCodecTest tests encoding, decoding and length calclulation for the given
// object.
func RunCodecDecodeTest(t *testing.T, c CodecTestCase) {
	t.Helper()

	t.Run(c.Name, func(t *testing.T) {
		t.Run("decode", func(t *testing.T) {
			// create a new instance of the same type as c.Struct
			typ := reflect.ValueOf(c.Struct).Type()
			var v reflect.Value
			switch typ.Kind() {
			case reflect.Ptr:
				v = reflect.New(typ.Elem()) // typ: *struct, v: *struct
			case reflect.Slice:
				v = reflect.New(typ) // typ: []x, v: *[]x
			default:
				t.Fatalf("%T is not a pointer or a slice", c.Struct)
			}

			if _, err := data.Decode(c.Bytes, v.Interface()); err != nil {
				t.Fatal(err)
			}

			// if v is a *[]x we need to dereference it before comparing it.
			if typ.Kind() == reflect.Slice {
				v = v.Elem()
			}

			klog.Infof("A:%v/n", c.Struct)
			//verify.Values(t, "", v.Interface(), c.Struct)
		})

		t.Run("encode", func(t *testing.T) {
			b, err := data.Encode(c.Struct)
			if err != nil {
				t.Fatal(err)
			}
			klog.Infof("A:[%d]%x/n", len(b), b)
			//			klog.Infof("B:[%d]%x/n",len(c.Bytes),c.Bytes)
			//			klog.Infof("Struct:%x/n",c.Bytes)
			///verify.Values(t, "", b, c.Bytes)
		})
	})

}

func RunCodecTest(t *testing.T, cases []CodecTestCase) {
	t.Helper()

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Run("decode", func(t *testing.T) {
				// create a new instance of the same type as c.Struct
				typ := reflect.ValueOf(c.Struct).Type()
				var v reflect.Value
				switch typ.Kind() {
				case reflect.Ptr:
					v = reflect.New(typ.Elem()) // typ: *struct, v: *struct
				case reflect.Slice:
					v = reflect.New(typ) // typ: []x, v: *[]x
				default:
					t.Fatalf("%T is not a pointer or a slice", c.Struct)
				}

				if _, err := data.Decode(c.Bytes, v.Interface()); err != nil {
					t.Fatal(err)
				}

				// if v is a *[]x we need to dereference it before comparing it.
				if typ.Kind() == reflect.Slice {
					v = v.Elem()
				}
				verify.Values(t, "", v.Interface(), c.Struct)
			})

			t.Run("encode", func(t *testing.T) {
				b, err := data.Encode(c.Struct)
				if err != nil {
					t.Fatal(err)
				}
				klog.Infof("A:%x/n", b)
				klog.Infof("B:%x/n", c.Bytes)
				verify.Values(t, "", b, c.Bytes)
			})
		})
	}
}

// RunCodecTest tests encoding, decoding and length calclulation for the given
// object.
func RunCodecHearderTest(t *testing.T, c CodecTestCase) {
	t.Helper()

	t.Run(c.Name, func(t *testing.T) {
		t.Run("decode", func(t *testing.T) {
			// create a new instance of the same type as c.Struct
			var st FireMessage
			typ := reflect.ValueOf(&st).Type()
			var v reflect.Value
			switch typ.Kind() {
			case reflect.Ptr:
				v = reflect.New(typ.Elem()) // typ: *struct, v: *struct
			case reflect.Slice:
				v = reflect.New(typ) // typ: []x, v: *[]x
			default:
				t.Fatalf("%T is not a pointer or a slice", c.Struct)
			}
			//data.NewFireBuffer(c.Bytes).Bytes()
			if _, err := data.Decode(c.Bytes, v.Interface()); err != nil {
				t.Fatal(err)
			}

			// if v is a *[]x we need to dereference it before comparing it.
			if typ.Kind() == reflect.Slice {
				v = v.Elem()
			}

			klog.Infof("Struct:%v/n", v.Interface().(*FireMessage))
			//verify.Values(t, "", v.Interface(), c.Struct)
		})
		//
		//t.Run("encode", func(t *testing.T) {
		//	b, err := data.Encode(c.Struct)
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//	klog.Infof("Hearder:[%d]%x/n",len(b),b)
		//	//			klog.Infof("B:[%d]%x/n",len(c.Bytes),c.Bytes)
		//	//			klog.Infof("Struct:%x/n",c.Bytes)
		//	///verify.Values(t, "", b, c.Bytes)
		//})
	})

}
