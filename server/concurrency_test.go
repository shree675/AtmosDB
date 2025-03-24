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

	go readWrite("0", wg, test)
	go readWrite("1", wg, test)

	wg.Wait()
}

func readWrite(id string, wg *sync.WaitGroup, test *testing.T) {
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
