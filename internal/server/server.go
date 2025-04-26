package server

import (
	"fmt"
	"net/http"

	"github.com/cmessinides/mnemonic/internal/bookmark"
	"github.com/cmessinides/mnemonic/internal/config"
	"github.com/labstack/echo/v4"
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
	e.HideBanner = true
	e.Debug = conf.Dev
	e.HTTPErrorHandler = customHTTPErrorHandler

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
	api.GET("/bookmarks/:id", b.Read)
	api.PATCH("/bookmarks/:id", b.Update)
	api.DELETE("/bookmarks/:id", b.Delete)

	return &Server{
		config: conf,
		e:      e,
	}
}

func (s *Server) Start() {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.e.Logger.Fatal(s.e.Start(address))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	c.Echo().DefaultHTTPErrorHandler(err, c)
}
