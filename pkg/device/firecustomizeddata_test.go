package device

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataInt64(t *testing.T) {
	assert.Equal(t, "123.44", parity(123.44, "string", "string"))
	assert.Equal(t, "123.44", parity(123.44, "int", "string"))
	assert.Equal(t, "4", parity(byte(0x04), "int", "string"))
	assert.Equal(t, "123", parity("123", "int", "string"))
	assert.Equal(t, "123.44", parity(123.44, "string", "string"))
	assert.Equal(t, "123.44", parity(float32(123.44), "float", "string"))
	assert.Equal(t, "123.44", parity(float64(123.44), "float", "string"))
	//assert.Equal(t, float64(1234567.44),parity(float32(1234567.44),"double"))
	assert.Equal(t, "123456789.44", parity(float64(123456789.44), "double", "string"))
	assert.Equal(t, "123456789", parity(int64(123456789), "double", "string"))
	assert.Equal(t, "123456789", parity(float64(123456789), "int", "string"))
}

func TestGetDeviceType(t *testing.T) {
	fmt.Println(getPointByDevicePoint("14.34.1.4.1.3.2.115.1.1"))
	assert.Equal(t, "2.115.1.1", getPointByDevicePoint("14.34.1.4.1.3.2.115.1.1"))
}
