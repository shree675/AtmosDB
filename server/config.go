package server

import (
	"log"
	"os"
)

var (
	DBVersion = ""
)

func initDBConfig() {
	version, err := os.ReadFile("./VERSION")
	if err != nil {
		log.Fatal("Failed to determine DB version:" + err.Error())
	}

	DBVersion = string(version)
}
