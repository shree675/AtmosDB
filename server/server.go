package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	st "atmosdb/server/store"
	t "atmosdb/types"
	"atmosdb/util"
)

func setValue(w http.ResponseWriter, r *http.Request) {
	ip, ok := parseInputPayload(w, r)
	if !ok {
		return
	}

	switch ip.Op {
	case int8(util.PUT):
		if st.StoreValue(&ip) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case int8(util.DELETE):
		st.DeleteValue(ip.Key)
		w.WriteHeader(http.StatusOK)
	case int8(util.DELTA):
		if st.UpdateValue(&ip) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Unsupported operation '"+strconv.Itoa(int(ip.Op))+"'", http.StatusBadRequest)
	}
}

func getValue(w http.ResponseWriter, r *http.Request) {
	ip, ok := parseInputPayload(w, r)
	if !ok {
		return
	}

	val, ok := st.GetValue(&ip)
	if !ok {
		http.Error(w, "No such key exists", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(val)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DBVersion)
}

func parseInputPayload(w http.ResponseWriter, r *http.Request) (t.InputPayload, bool) {
	if r.Method != "POST" {
		http.Error(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
		return t.InputPayload{}, false
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error while reading request body: "+err.Error(), http.StatusInternalServerError)
		return t.InputPayload{}, false
	}

	var ip t.InputPayload
	if err := json.Unmarshal(body, &ip); err != nil {
		http.Error(w, "Error while parsing request body: "+err.Error(), http.StatusBadRequest)
		return t.InputPayload{}, false
	}

	return ip, true
}

func Start() {
	http.HandleFunc("/set", setValue)
	http.HandleFunc("/get", getValue)
	http.HandleFunc("/version", getVersion)

	// explicitly initializing instead of using init() for code readability
	st.Init()
	st.InitScheduler()

	fmt.Println("AtmosDB " + DBVersion)
	log.Println("Initialized workers")
	log.Println("Running AtmosDB at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
