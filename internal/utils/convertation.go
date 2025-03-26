package utils

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

func ConvertToStringSlice(val interface{}) []string {
	items := val.([]interface{})
	result := make([]string, len(items))
	for i, v := range items {
		result[i] = v.(string)
	}
	return result
}

func ConvertToIntMap(val interface{}) map[string]int {
	result := make(map[string]int)

	if val == nil {
		return result
	}

	rawMap, ok := val.(map[interface{}]interface{})
	if !ok {
		return result
	}

	for k, v := range rawMap {
		key := fmt.Sprintf("%v", k)

		switch num := v.(type) {
		case int64:
			result[key] = int(num)
		case int32:
			result[key] = int(num)
		case int16:
			result[key] = int(num)
		case int8:
			result[key] = int(num)
		case int:
			result[key] = num
		case uint64:
			if num > math.MaxInt {
				result[key] = math.MaxInt
			} else {
				result[key] = int(num)
			}
		case uint32:
			result[key] = int(num)
		case uint16:
			result[key] = int(num)
		case uint8:
			result[key] = int(num)
		case uint:
			if num > math.MaxInt {
				result[key] = math.MaxInt
			} else {
				result[key] = int(num)
			}
		case float64:
			result[key] = int(num)
		case float32:
			result[key] = int(num)
		default:
			result[key] = 0
		}
	}

	return result
}

func ConvertToBoolMap(val interface{}) map[string]bool {
	m := val.(map[interface{}]interface{})
	result := make(map[string]bool)
	for k, v := range m {
		result[k.(string)] = v.(bool)
	}
	return result
}

func InterfaceToUint64(val interface{}) uint64 {
	switch v := val.(type) {
	case uint64:
		return v
	case int64, int32, int16, int8, int:
		return uint64(reflect.ValueOf(val).Int())
	case uint, uint32, uint16, uint8:
		return uint64(reflect.ValueOf(val).Uint())
	default:
		return 0
	}
}

func InterfaceToInt64(val interface{}) (int64, error) {
	switch v := val.(type) {
	case int64:
		return v, nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, errors.New("value too large for int64")
		}
		return int64(v), nil
	case int, int8, int16, int32:
		return reflect.ValueOf(val).Int(), nil
	case uint, uint8, uint16, uint32:
		return int64(reflect.ValueOf(val).Uint()), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", val)
	}
}
