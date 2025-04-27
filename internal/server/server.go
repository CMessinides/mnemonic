package server

import (
	"fmt"

	"github.com/cmessinides/mnemonic/internal/bookmark"
	"github.com/cmessinides/mnemonic/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"method":"${method}","uri":"${uri}","status":"${status}"}` + "\n",
	}))
	e.HTTPErrorHandler = customHTTPErrorHandler

	assets := NewAssetsFS(AssetConfig{
		PublicPath: "/assets",
		Dev:        conf.Dev,
	})
	e.StaticFS(assets.PublicPath, assets)

	funcs := assets.TemplateFuncs()
	t := NewTemplate(TemplateConf{
		Dev:         conf.Dev,
		CustomFuncs: funcs,
	})
	e.Renderer = t

	h := &homeController{bookmarks: bookmarks}
	e.GET("/", h.Show)

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
