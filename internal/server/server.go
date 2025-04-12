package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo
}

func NewServer() *Server {
	e := echo.New()

	t := NewTemplate()
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", map[string]interface{}{"Name": "Cam"})
	})

	e.StaticFS("/assets", NewAssetsFS())

	return &Server{
		e: e,
	}
}

func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(":3000"))
}
