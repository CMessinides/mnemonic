package server

import (
	"fmt"
	"net/http"

	"github.com/cmessinides/mnemonic/internal/bookmark"
	"github.com/cmessinides/mnemonic/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Server struct {
	config *Config
	e      *echo.Echo
}

type Config struct {
	config.ServerConfig
	Dev       bool
	LookupEnv config.LookupEnv
}

func NewServer(conf *Config, bookmarks bookmark.BookmarkStore) *Server {
	e := echo.New()

	if conf.Dev {
		e.Logger.SetLevel(log.DEBUG)
		e.Logger.Info("dev mode enabled")
	}

	t := NewTemplate()
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", map[string]any{"Name": "Cam"})
	})

	e.StaticFS("/assets", GetAssetsFS(AssetConfig{
		Dev: conf.Dev,
	}))

	api := e.Group("/api/v1")

	b := &bookmarksAPI{store: bookmarks}
	api.GET("/bookmarks", b.List)
	api.POST("/bookmarks", b.Create)

	return &Server{
		config: conf,
		e:      e,
	}
}

func (s *Server) Start() {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.e.Logger.Fatal(s.e.Start(address))
}
