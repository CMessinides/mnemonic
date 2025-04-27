package server

import (
	"embed"
	"html/template"
	"io"
	"maps"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

//go:embed public/views/*.html
var views embed.FS

//go:embed public/views/layouts/*.html
var layouts embed.FS

//go:embed public/views/partials/*.html
var partials embed.FS

type TemplateConf struct {
	Dev         bool
	CustomFuncs template.FuncMap
}

type DevTemplate struct {
	funcs template.FuncMap
}

func (d *DevTemplate) ExecuteTemplate(w io.Writer, name string, data any) error {
	t := template.New("default.html").Funcs(d.funcs)
	t, err := t.ParseFS(os.DirFS("internal/server/public/views/layouts"), "*.html")
	if err != nil {
		return err
	}

	t, err = t.ParseFS(os.DirFS("internal/server/public/views/partials"), "*.html")
	if err != nil {
		return err
	}

	t, err = t.ParseFiles("internal/server/public/views/" + name)
	if err != nil {
		return err
	}

	return t.ExecuteTemplate(w, name, data)
}

type EmbeddedTemplate struct {
	shared *template.Template
}

func (e *EmbeddedTemplate) ExecuteTemplate(w io.Writer, name string, data any) error {
	shared, err := e.shared.Clone()
	if err != nil {
		return err
	}

	t, err := shared.ParseFS(views, "public/views/"+name)
	if err != nil {
		return err
	}

	return t.ExecuteTemplate(w, name, data)
}

func NewTemplate(conf TemplateConf) *echo.TemplateRenderer {
	f := template.FuncMap{
		"formatISOTimestamp": func(t time.Time) string {
			return t.Format(time.RFC3339)
		},
		"formatRelativeTime": func(t time.Time) string {
			return time.Since(t).String()
		},
	}

	maps.Copy(f, conf.CustomFuncs)

	if !conf.Dev {
		shared := template.New("default.html").Funcs(f)
		shared = template.Must(shared.ParseFS(layouts))
		shared = template.Must(shared.ParseFS(partials))
		return &echo.TemplateRenderer{
			Template: &EmbeddedTemplate{
				shared: shared,
			},
		}
	} else {
		return &echo.TemplateRenderer{
			Template: &DevTemplate{funcs: f},
		}
	}
}
