package types

import (
	"net/http"
)

type SessionConfig struct {
	Client *http.Client
	SId    string
	Conn   string
}
