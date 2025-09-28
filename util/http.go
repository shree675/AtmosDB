package util

import (
	t "atmosdb/types"
	"encoding/json"
	"io"
	"net/http"
)

func ParseInputPayload(w http.ResponseWriter, r *http.Request) (*t.InputPayload, bool) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error while reading request body: "+err.Error(), http.StatusInternalServerError)
		return nil, false
	}

	var ip t.InputPayload
	if err := json.Unmarshal(body, &ip); err != nil {
		http.Error(w, "Error while parsing request body: "+err.Error(), http.StatusBadRequest)
		return nil, false
	}

	return &ip, true
}

func ParseInputSubscriptionPayload(w http.ResponseWriter, r *http.Request) (*t.InputSubscriptionPayload, bool) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error while reading request body: "+err.Error(), http.StatusInternalServerError)
		return nil, false
	}

	var isp t.InputSubscriptionPayload
	if err := json.Unmarshal(body, &isp); err != nil {
		http.Error(w, "Error while parsing request body: "+err.Error(), http.StatusBadRequest)
		return nil, false
	}

	return &isp, true
}
