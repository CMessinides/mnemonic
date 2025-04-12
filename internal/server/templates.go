package server

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

//go:embed public/views/*.html
var templates embed.FS

type Template struct {
	templates *template.Template
}

func NewTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseFS(templates, "public/views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl = template.Must(tmpl.ParseFS(templates, "public/views/"+name))
	return tmpl.ExecuteTemplate(w, name, data)
}
