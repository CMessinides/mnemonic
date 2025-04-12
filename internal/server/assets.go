package server

import (
	"embed"
	"io/fs"

	"github.com/labstack/echo/v4"
)

//go:embed all:public/assets
var assets embed.FS

func NewAssetsFS() fs.FS {
	return echo.MustSubFS(assets, "public/assets")
}
