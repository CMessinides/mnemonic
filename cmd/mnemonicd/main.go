package main

import (
	"os"

	"github.com/cmessinides/mnemonic/internal/server"
)

var DevMode = "off"

func main() {
	s := server.NewServer(server.ServerConfig{
		// go run -ldflags "-X main.DevMode=on" ./cmd/mnemonicd
		Dev:    DevMode == "on",
		GetEnv: os.LookupEnv,
	})
	s.Start()
}
