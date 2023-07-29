package utils

import (
	"reflect"
)

func ZeroArray(bytes *[]byte) {
	for i := 0; i < len(*bytes); i++ {
		(*bytes)[i] = 0x00
	}
}

func MaxArray(values []uint32) (value uint32, index int) {
	biggestValue := uint32(0)
	biggestIndex := 0

	for actualIndex, actualValue := range values {
		if actualValue > biggestValue {
			biggestValue = actualValue
			biggestIndex = actualIndex
		}
	}

	return biggestValue, biggestIndex
}

func MaxUint32(a uint32, b uint32) uint32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func ClearMap(mapToClearGeneric interface{}) {
	mapToClear := reflect.ValueOf(mapToClearGeneric)

	if mapToClear.Kind() != reflect.Map {
		panic("WTF bro?!")
		return
	}

	for _, key := range mapToClear.MapKeys() {
		mapToClear.SetMapIndex(key, reflect.Value{})
	}
}
