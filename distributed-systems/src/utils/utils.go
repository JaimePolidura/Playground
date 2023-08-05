package utils

import (
	"reflect"
	"sync"
)

func ZeroArray(bytes *[]byte) {
	for i := 0; i < len(*bytes); i++ {
		(*bytes)[i] = 0x00
	}
}

func GetInt32FromSyncMap(syncMap *sync.Map, key any) int32 {
	value, contained := syncMap.Load(key)

	if contained {
		valueInt, _ := value.(int32)
		return valueInt
	} else {
		return 0
	}
}

func MaxUint32(a uint32, b uint32) uint32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinInt32(a int32, b int32) int32 {
	if a > b {
		return b
	} else {
		return a
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
