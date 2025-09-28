package store

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	t "atmosdb/types"
	"atmosdb/util"
)

var stOut map[string][]*http.ResponseWriter
var stream chan t.StreamItem
var stSem map[int]*sync.RWMutex

func InitStreamReader() {
	stOut = make(map[string][]*http.ResponseWriter)
	stream = make(chan t.StreamItem, util.StreamQ)
	stSem = make(map[int]*sync.RWMutex, util.Concurrency)
	for i := 0; i < util.Concurrency; i++ {
		stSem[i] = &sync.RWMutex{}
	}

	go func() {
		for {
			it := <-stream
			hash := hashCode(it.Key)

			stSem[hash].RLock()
			for _, w := range stOut[it.Key] {
				_, err := fmt.Fprintf(*w, "%v", it.Val)
				if err != nil {
					log.Println("[STREAM_READER] failed to write to writer", err)
				} else {
					(*w).(http.Flusher).Flush()
				}
			}
			stSem[hash].RUnlock()

			if it.Val == util.StreamDeleteId {
				deleteStreamKey(it.Key)
			}
		}
	}()
}

func AddWriter(isp *t.InputSubscriptionPayload, w *http.ResponseWriter) {
	hash := hashCode(isp.Key)

	stSem[hash].Lock()

	_, exists := stOut[isp.Key]
	if !exists {
		stOut[isp.Key] = make([]*http.ResponseWriter, 0)
	}
	stOut[isp.Key] = append(stOut[isp.Key], w)

	stSem[hash].Unlock()

	log.Println("[" + isp.SId + "] Added new subscriber for key: " + isp.Key)
}

func RemoveWriter(isp *t.InputSubscriptionPayload, w *http.ResponseWriter) {
	hash := hashCode(isp.Key)

	stSem[hash].Lock()
	defer stSem[hash].Unlock()

	for i, val := range stOut[isp.Key] {
		if val == w {
			stOut[isp.Key] = append(stOut[isp.Key][:i], stOut[isp.Key][i+1:]...)
			break
		}
	}
}

func deleteStreamKey(key string) {
	hash := hashCode(key)

	stSem[hash].Lock()
	defer stSem[hash].Unlock()

	delete(stOut, key)
}
