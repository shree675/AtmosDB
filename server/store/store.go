package store

import (
	"hash/fnv"
	"log"
	"sync"
	"time"

	"github.com/dgrijalva/lfu-go"

	t "atmosdb/types"
	"atmosdb/util"
)

/*
use normal map with custom locking implementation
instead of sync.Map to reduce memory footprint,
evident in case of large number of key-value pairs
*/
var data map[string]t.Data
var sem map[int]*sync.RWMutex
var cache lfu.Cache

func Init() {
	data = make(map[string]t.Data)

	cache = *lfu.New()
	cache.UpperBound = util.CacheUB
	cache.LowerBound = util.CacheUB - 1

	sem = make(map[int]*sync.RWMutex, util.Concurrency)
	for i := 0; i < util.Concurrency; i++ {
		sem[i] = &sync.RWMutex{}
	}
}

func StoreValue(ip *t.InputPayload) bool {
	hash := hashCode(ip.Key)

	sem[hash].Lock()

	ttl := ip.Ttl
	now := time.Now().UnixMilli()
	dt := util.Types[ip.Type]
	val, ok := dt.Convert(ip.Val)
	if !ok {
		return false
	}
	data[ip.Key] = t.Data{
		Key:    ip.Key,
		Ttl:    ttl,
		Val:    val,
		OpTime: now,
		Type:   ip.Type,
	}

	sem[hash].Unlock()

	if ttl != 0 {
		go AddExpiryItem(ip.Key, now+int64(ttl)*1000)
	}

	return true
}

func UpdateValue(ip *t.InputPayload) bool {
	hash := hashCode(ip.Key)

	sem[hash].Lock()
	defer sem[hash].Unlock()

	edat, exists := data[ip.Key]
	if !exists {
		log.Println(util.GetYellowStr("[ERROR] [" + ip.SId + "] No such key exists: " + ip.Key))
		return false
	} else if edat.Type != int8(util.INT) {
		log.Println(util.GetYellowStr("[ERROR] [" + ip.SId + "] Wrong type for key: " + ip.Key + ", expected (int32)"))
		return false
	}

	dt := util.Types[ip.Type]
	delta, _ := dt.Convert(ip.Val)
	edat.Val = edat.Val.(int) + delta.(int)
	data[ip.Key] = edat

	return true
}

func GetValue(ip *t.InputPayload) (t.OutputPayload, bool) {
	hash := hashCode(ip.Key)

	sem[hash].RLock()
	defer sem[hash].RUnlock()

	d, ok := data[ip.Key]
	if !ok {
		log.Println("["+ip.SId+"] No such key exists:", ip.Key)
		return t.OutputPayload{}, false
	}

	return t.OutputPayload{Key: d.Key, Val: d.Val, Type: d.Type}, true
}

func DeleteValue(key string) {
	hash := hashCode(key)

	sem[hash].Lock()
	defer sem[hash].Unlock()

	delete(data, key)
}

func hashCode(key string) int {
	if cache.Get(key) != nil {
		return cache.Get(key).(int)
	}

	algorithm := fnv.New32a()
	algorithm.Write([]byte(key))
	v := int(algorithm.Sum32() % uint32(util.Concurrency))

	cache.Set(key, v)
	return v
}
