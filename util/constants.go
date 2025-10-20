package util

import (
	"reflect"

	t "atmosdb/types"
)

var Types = map[int8]t.DataType{
	int8(INT):    t.IntType{Type: reflect.TypeOf(int32(0))},
	int8(FLOAT):  t.FloatType{Type: reflect.TypeOf(float32(0))},
	int8(STRING): t.StringType{Type: reflect.TypeOf("")},
}

const Concurrency int = 16
const TtlFreq int64 = 2e9
const CacheUB int = 1000
const StreamQ int = 1000
const StreamDeleteId string = "<nil>"
const StreamBufSize int = 1024
