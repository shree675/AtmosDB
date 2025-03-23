package types

type InputPayload struct {
	SId  string
	Key  string
	Val  string `json:"value"`
	Ttl  uint32
	Type int8 // see util/constants.go
	Op   int8 // see util/constants.go
}

type OutputPayload struct {
	Key  string
	Val  any
	Type int8
}
