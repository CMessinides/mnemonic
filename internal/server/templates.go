package server

import (
	"embed"
	"html/template"
	"io"
	"maps"
	"time"

	"github.com/labstack/echo/v4"
)

//go:embed public/views/*.html
var templates embed.FS

type Template struct {
	templates *template.Template
}

func NewTemplate(funcs template.FuncMap) *Template {
	f := template.FuncMap{
		"formatISOTimestamp": func(t time.Time) string {
			return t.Format(time.RFC3339)
		},
		"formatRelativeTime": func(t time.Time) string {
			return time.Since(t).String()
		},
	}

	maps.Copy(f, funcs)

	return &Template{
		templates: template.Must(
			template.
				New("mnemonic").
				Funcs(f).
				ParseFS(templates, "public/views/*.html"),
		),
	}
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl = template.Must(tmpl.ParseFS(templates, "public/views/"+name))
	return tmpl.ExecuteTemplate(w, name, data)
}
