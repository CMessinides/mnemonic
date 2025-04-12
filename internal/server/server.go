package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Server struct {
	e *echo.Echo
}

type ServerConfig struct {
	Dev    bool
	GetEnv func(string) (string, bool)
}

func NewServer(conf ServerConfig) *Server {
	e := echo.New()

	if conf.Dev {
		e.Logger.SetLevel(log.DEBUG)
		e.Logger.Info("dev mode enabled")
	}

	t := NewTemplate()
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", map[string]interface{}{"Name": "Cam"})
	})

	e.StaticFS("/assets", GetAssetsFS(AssetConfig{
		Dev: conf.Dev,
	}))

	return &Server{
		e: e,
	}
}

func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(":3000"))
}
