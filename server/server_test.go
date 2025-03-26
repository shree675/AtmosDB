package server

import (
	"strconv"
	"sync"
	"testing"
	"time"

	st "atmosdb/server/store"
	t "atmosdb/types"
	"atmosdb/util"
)

func TestConcurrentReadsWrites(test *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	st.Init()

	go rwWorker("0", wg, test)
	go rwWorker("1", wg, test)

	wg.Wait()
}

func TestTTLPriority(test *testing.T) {
	st.Init()
	st.InitScheduler()

	ip1 := t.InputPayload{
		SId:  "0",
		Key:  "k1",
		Val:  "1",
		Ttl:  4,
		Type: int8(util.INT),
		Op:   int8(util.PUT),
	}
	ip2 := t.InputPayload{
		SId:  "0",
		Key:  "k2",
		Val:  "2",
		Ttl:  1,
		Type: int8(util.INT),
		Op:   int8(util.PUT),
	}

	st.StoreValue(&ip1)
	st.StoreValue(&ip2)

	qp1 := t.InputPayload{
		SId: "0",
		Key: "k1",
	}
	qp2 := t.InputPayload{
		SId: "0",
		Key: "k2",
	}

	time.Sleep(3 * time.Second)

	val, gtOk1 := st.GetValue(&qp1)
	if !gtOk1 {
		test.Error("Failed to get value for 'k1'")
	}
	if val.Val != 1 {
		test.Errorf("Expected 1, got %d\n", val.Val)
	}

	_, gtOk2 := st.GetValue(&qp2)
	if gtOk2 {
		test.Error("Key 'k2' not expired")
	}

	time.Sleep(5 * time.Second)

	_, ok := st.GetValue(&qp1)
	if ok {
		test.Error("Key 'k1' not expired")
	}
}

func rwWorker(id string, wg *sync.WaitGroup, test *testing.T) {
	defer wg.Done()

	for i := 0; i < 10; i++ {
		ip1 := t.InputPayload{
			SId:  id,
			Key:  "k1",
			Val:  strconv.Itoa(i),
			Type: int8(util.INT),
			Op:   int8(util.PUT),
		}

		stOk := st.StoreValue(&ip1)
		if !stOk {
			test.Error("Failed to store value")
		}

		time.Sleep(100 * time.Millisecond)

		ip2 := t.InputPayload{
			SId: id,
			Key: "k1",
		}

		_, gtOk := st.GetValue(&ip2)
		if !gtOk {
			test.Error("Failed to get value")
		}
	}
}
