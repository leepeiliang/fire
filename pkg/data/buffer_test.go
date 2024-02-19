package data

import "testing"

func WriteByteStringTest(t *testing.T, cases []CodecTestCase) {
	buf := NewBuffer(nil)
	buf.WriteByteString(nil)
}
