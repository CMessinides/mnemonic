package server

import (
	"embed"
	"io/fs"
	"os"

	"github.com/labstack/echo/v4"
)

//go:embed all:public/assets
var assets embed.FS

func NewAssetsFS() fs.FS {
	return echo.MustSubFS(assets, "public/assets")
}

type AssetConfig struct {
	Dev bool
}

func GetAssetsFS(conf AssetConfig) fs.FS {
	if conf.Dev {
		return os.DirFS("internal/server/public/assets")
	} else {
		return echo.MustSubFS(assets, "public/assets")
	}
}
