package server

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"

	t "atmosdb/types"
)

var (
	DBVersion = ""
)

func init() {
	_, configPath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to load DB config")
	}
	configPath = filepath.Dir(filepath.Dir(configPath))

	config, err := os.ReadFile(filepath.Join(configPath, "dbconfig.json"))
	if err != nil {
		log.Fatal("Failed to load DB config: " + err.Error())
	}

	sc := readConfig(config)
	setVersion(sc)
}

func readConfig(config []byte) t.ServerConfig {
	var sc t.ServerConfig
	if err := json.Unmarshal(config, &sc); err != nil {
		log.Fatal("Failed to load DB config: " + err.Error())
	}

	return sc
}

func setVersion(sc t.ServerConfig) {
	DBVersion = sc.Version

	if DBVersion == "" {
		log.Fatal("DB version missing in config")
	}
}
