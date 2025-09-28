package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	st "atmosdb/server/store"
	"atmosdb/util"
)

func setValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip, ok := util.ParseInputPayload(w, r)
	if !ok {
		return
	}

	switch ip.Op {
	case int8(util.PUT):
		if st.StoreValue(ip) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case int8(util.DELETE):
		st.DeleteValue(ip.Key)
		w.WriteHeader(http.StatusOK)
	case int8(util.DELTA):
		if st.UpdateValue(ip) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Unsupported operation '"+strconv.Itoa(int(ip.Op))+"'", http.StatusBadRequest)
	}
}

func getValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip, ok := util.ParseInputPayload(w, r)
	if !ok {
		return
	}

	val, ok := st.GetValue(ip)
	if !ok {
		http.Error(w, "No such key exists", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(val)
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
		return
	}

	isp, ok := util.ParseInputSubscriptionPayload(w, r)
	if !ok {
		return
	}

	exists := st.HasKey(isp.Key)
	if !exists {
		http.Error(w, "No such key exists", http.StatusNotFound)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	done := r.Context().Done()

	st.AddWriter(isp, &w)

	<-done
	st.RemoveWriter(isp, &w)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DBVersion)
}

func Start() {
	http.HandleFunc("/set", setValue)
	http.HandleFunc("/get", getValue)
	http.HandleFunc("/subscribe", subscribe)
	http.HandleFunc("/version", getVersion)

	fmt.Println("AtmosDB " + DBVersion)

	// explicitly initializing instead of using init() for code readability
	st.Init()
	st.InitScheduler()
	st.InitStreamReader()

	log.Println("Initialized workers")
	log.Println("Running AtmosDB at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
