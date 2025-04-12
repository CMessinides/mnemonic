#!/bin/sh
go run -ldflags "-X main.DevMode=on" ./cmd/mnemonicd
