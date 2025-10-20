package types

import (
	"net/http"
)

type ServerConfig struct {
	Version string
}

type SessionConfig struct {
	Client    *http.Client
	SSEClient *http.Client
	SId       string
	Conn      string
}
