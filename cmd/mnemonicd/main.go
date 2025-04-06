package main

import "github.com/cmessinides/mnemonic/internal/server"

func main() {
	s := server.NewServer()
	s.Start()
}
