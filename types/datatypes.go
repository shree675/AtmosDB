package types

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

type DataType interface {
	Convert(val string) (any, bool)
	GetType() reflect.Type
}

type IntType struct {
	Type reflect.Type
}

func (it IntType) Convert(val string) (any, bool) {
	v, err := strconv.Atoi(val)
	ok := true
	if err != nil {
		ok = false
		log.Println("[ERROR] Cannot convert " + val + " to int type")
	}
	return v, ok
}

func (it IntType) GetType() reflect.Type {
	return it.Type
}

type FloatType struct {
	Type reflect.Type
}

func (it FloatType) Convert(val string) (any, bool) {
	v, err := strconv.ParseFloat(val, 32)
	ok := true
	if err != nil {
		ok = false
		log.Println("[ERROR] Cannot convert " + val + " to int type")
	}
	return float32(v), ok
}

func (it FloatType) GetType() reflect.Type {
	return it.Type
}

type StringType struct {
	Type reflect.Type
}

func (it StringType) Convert(val string) (any, bool) {
	return strings.Trim(val, "\""), true
}

func (it StringType) GetType() reflect.Type {
	return it.Type
}
