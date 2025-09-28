package types

import (
	"github.com/Workiva/go-datastructures/queue"
)

type Data struct {
	Key    string
	Ttl    uint32 // seconds
	OpTime int64
	Val    any
	Type   int8
}

type ExpiryItem struct {
	Key  string
	Time int64
}

func (it ExpiryItem) Compare(other queue.Item) int {
	if it.Time > other.(ExpiryItem).Time {
		return 1
	} else if it.Time < other.(ExpiryItem).Time {
		return -1
	} else {
		return 0
	}
}

type StreamItem struct {
	Key string
	Val any
}
