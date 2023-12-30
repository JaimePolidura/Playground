package interpreter

import (
	"errors"
	"reflect"
	"strconv"
)

func castBoolean(value any) (bool, error) {
	switch value.(type) {
	case float64:
		return value.(float64) > 0, nil
	case bool:
		return value.(bool), nil
	case nil:
		return false, nil
	default:
		return false, errors.New("Cannot take " + reflect.TypeOf(value).Name() + " as boolean")
	}
}

func castString(value any) (string, error) {
	switch value.(type) {
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64), nil
	case bool:
		if value.(bool) {
			return "true", nil
		} else {
			return "false", nil
		}
	case nil:
		return "", nil
	case string:
		return value.(string), nil
	case LoxInstance:
		return value.(LoxInstance).KClass.Name, nil
	default:
		return "", errors.New("Cannot take " + reflect.TypeOf(value).Name() + " as number")
	}
}

func castNumber(value any) (float64, error) {
	switch value.(type) {
	case float64:
		return value.(float64), nil
	case bool:
		if value.(bool) {
			return 1, nil
		} else {
			return 0, nil
		}
	case nil:
		return 0, nil
	default:
		return -1, errors.New("Cannot take " + reflect.TypeOf(value).Name() + " as number")
	}
}
