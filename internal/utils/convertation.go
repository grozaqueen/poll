package utils

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
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

func ToInt64(val interface{}) (int64, error) {
	if val == nil {
		return 0, errors.New("nil value")
	}

	switch v := val.(type) {
	case int64:
		return v, nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, errors.New("value too large for int64")
		}
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case float64:
		if v > math.MaxInt64 || v < math.MinInt64 {
			return 0, errors.New("value out of int64 range")
		}
		return int64(v), nil
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse string as int64: %v", err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported type for int64: %T", val)
	}
}

func ToString(val interface{}) string {
	if s, ok := val.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", val)
}
