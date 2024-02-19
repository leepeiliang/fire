package device

import (
	"github.com/shopspring/decimal"
	"k8s.io/klog/v2"
	"reflect"

	"math"
	"strconv"
)

func parity(data interface{}, propertyDataType, outDataType string) interface{} {
	switch outDataType {
	case "string":
		return data2string(data, propertyDataType)
	case "double":
		return data2float64(data)
	case "float":
		return data2float32(data)
	case "boolean":
		return data2bool(data)
	case "bytes":
		return data.(byte)
	case "int":
		return data2int(data)
	default:
		return data2string(data, propertyDataType)
	}
	return data2string(data, propertyDataType)
}
func data2int(data interface{}) int64 {
	switch data.(type) {
	case string:
		f, err := strconv.ParseInt(data.(string), 10, 64)
		if err != nil {
			return 0
		}
		if f < math.MaxInt64 {
			return f
		}
		return math.MaxInt64
	case uint, uint64, uint32, uint16, int64, int, int32, int16:
		return data.(int64)
	case float32:
		if data.(float32) < math.MaxInt64 {
			return int64(math.Ceil(float64(data.(float32))))
		}
		return math.MaxInt64
	case float64:
		if data.(float64) < math.MaxInt64 {
			return int64(math.Ceil(data.(float64)))
		}
		return math.MaxInt64
	case byte:
		return int64(data.(byte))
	case bool:
		if data.(bool) == true {
			return 1
		}
		return 0
	default:
		return 0
	}
	return 0
}

func data2string(data interface{}, propertyDataType string) string {

	klog.V(4).Infof("dataType:%+v+v", reflect.ValueOf(data), reflect.TypeOf(data))
	switch data.(type) {
	case string:
		if propertyDataType == "float" {
			decimalNum, err := decimal.NewFromString(data.(string))
			if err != nil {
				return "0.00"
			}
			return decimalNum.StringFixed(2)
		}
		if propertyDataType == "double" {
			decimalNum, err := decimal.NewFromString(data.(string))
			if err != nil {
				return "0.00"
			}
			return decimalNum.StringFixed(2)
		}
		return data.(string)
	case uint:
		return strconv.FormatInt(int64(data.(uint)), 10)
	case uint64:
		return strconv.FormatInt(int64(data.(uint64)), 10)
	case uint16:
		return strconv.FormatInt(int64(data.(uint16)), 10)
	case int64:
		return strconv.FormatInt(data.(int64), 10)
	case int:
		return strconv.FormatInt(int64(data.(int)), 10)
	case int32:
		return strconv.FormatInt(int64(data.(int32)), 10)
	case uint32:
		return strconv.FormatInt(int64(data.(uint32)), 10)
	case float32:
		return strconv.FormatFloat(float64(data.(float32)), 'f', 2, 32)
	case float64:
		return strconv.FormatFloat(data.(float64), 'f', 2, 64)
	case byte:
		return strconv.Itoa(int(data.(byte)))
	case bool:
		return strconv.FormatBool(data.(bool))
	default:
		return "N/A"
	}
	return "N/A"
}

func data2float64(data interface{}) float64 {
	switch data.(type) {
	case string:
		f, err := strconv.ParseFloat(data.(string), 64)
		if err != nil {
			return 0.00
		}
		if f < math.MaxFloat64 {
			return f
		}
		return math.MaxFloat64
	case float32:
		return float64(data.(float32))
	case float64:
		return data.(float64)
	case uint, uint64, uint32, uint16:
		f, err := strconv.ParseFloat(strconv.FormatUint(data.(uint64), 10), 64)
		if err != nil {
			return 0.00
		}
		return f
	case int64, int, int32, int16:
		f, err := strconv.ParseFloat(strconv.FormatInt(data.(int64), 10), 64)
		if err != nil {
			return 0.00
		}
		return f
	default:
		return 0.00
	}

	return 0.00
}

// 如果转换的数据大于math.MaxFloat32，那么返回math.MaxFloat32
func data2float32(data interface{}) float32 {
	switch data.(type) {
	case string:
		f, err := strconv.ParseFloat(data.(string), 32)
		if err != nil {
			return 0.00
		}
		if f < math.MaxFloat32 {
			return float32(f)
		}
		return math.MaxFloat32
	case float32:
		return data.(float32)
	case float64:
		if data.(float64) < math.MaxFloat32 {
			return float32(data.(float64))
		}
		return math.MaxFloat32
	case uint, uint64, uint32, uint16:
		f, err := strconv.ParseFloat(strconv.FormatUint(data.(uint64), 10), 32)
		if err != nil {
			return 0.00
		}
		if f < math.MaxFloat32 {
			return float32(f)
		}
		return math.MaxFloat32
	case int64, int, int32, int16:
		f, err := strconv.ParseFloat(strconv.FormatInt(data.(int64), 10), 64)
		if err != nil {
			return 0.00
		}
		if f < math.MaxFloat32 {
			return float32(f)
		}
		return math.MaxFloat32
	default:
		return 0.00
	}

	return 0.00
}

func data2bool(data interface{}) bool {
	switch data.(type) {
	case string:
		if data.(string) == "true" {
			return true
		}
		return false
	case float32:
		return data.(float32) > 0
	case float64:
		return data.(float64) > 0
	case uint, uint64, uint32, uint16:
		return data.(uint) > 0
	case int64, int, int32, int16:
		return data.(int64) > 0
	case bool:
		return data.(bool)
	case byte:
		return data.(byte) > 0
	default:
		return false
	}

	return false
}
