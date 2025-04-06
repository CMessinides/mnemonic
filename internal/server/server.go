package server

import (
	"embed"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed public/views
var templates embed.FS

type Server struct {
	e *echo.Echo
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewServer() *Server {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseFS(templates, "public/views/*.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "hello", "Cam")
	})

	return &Server{
		e: e,
	}
}

func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(":3000"))
}
