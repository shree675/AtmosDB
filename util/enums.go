package util

type BaseType int8
type OpType int8
type Command string

const (
	INT BaseType = iota
	FLOAT
	STRING
)

const (
	PUT OpType = iota
	DELETE
	DELTA
)

const (
	GET          Command = "GET"
	SETINT       Command = "SETINT"
	SETFLOAT     Command = "SETFLOAT"
	SETSTR       Command = "SETSTR"
	SETINT_TTL   Command = "SETINT.TTL"
	SETFLOAT_TTL Command = "SETFLOAT.TTL"
	SETSTR_TTL   Command = "SETSTR.TTL"
	DEL          Command = "DEL"
	INCR         Command = "INCR"
	DECR         Command = "DECR"
	EXISTS       Command = "EXISTS"
)
