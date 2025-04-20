package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/cmessinides/mnemonic/internal/bookmark"
	"github.com/cmessinides/mnemonic/internal/config"
	"github.com/cmessinides/mnemonic/internal/server"
	_ "modernc.org/sqlite"
)

var DevMode = "off"

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

func main() {
	conf, err := config.ReadConfig(
		xdg.Home,
		xdg.ConfigHome,
		xdg.DataHome,
		os.LookupEnv,
		fileExists,
		os.ReadFile,
	)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("sqlite", conf.DBConnString())
	if err != nil {
		log.Fatalln(err)
	}

	bookmarks := bookmark.NewSQLiteBookmarkStore(db)
	err = bookmarks.Init()
	if err != nil {
		log.Fatalln(err)
	}

	s := server.NewServer(&server.Config{
		// go run -ldflags "-X main.DevMode=on" ./cmd/mnemonicd
		ServerConfig: *conf.Server,
		Dev:          DevMode == "on",
		LookupEnv:    os.LookupEnv,
	}, bookmarks)

	s.Start()
}
