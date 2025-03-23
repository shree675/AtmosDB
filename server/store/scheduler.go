package store

import (
	"log"
	"time"

	"github.com/Workiva/go-datastructures/queue"

	t "atmosdb/types"
	"atmosdb/util"
)

var pq queue.PriorityQueue // pq is thread-safe

func InitScheduler() {
	pq = *queue.NewPriorityQueue(100, true)

	go expireEntries()
}

func AddExpiryItem(key string, time int64) {
	pq.Put(t.ExpiryItem{Key: key, Time: time})
}

func expireEntries() {
	ticker := time.NewTicker(time.Duration(util.TtlFreq))

	for range ticker.C {
		if !pq.Empty() {
			v := pq.Peek().(t.ExpiryItem)
			if v.Time <= time.Now().UnixMilli() {
				pq.Get(1)
				DeleteValue(v.Key)
				log.Println("Expired key " + v.Key)
			}
		}
	}
}
