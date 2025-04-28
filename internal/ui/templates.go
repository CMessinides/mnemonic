package ui

import (
	"embed"
	"html/template"
	"io"
	"maps"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

//go:embed views/*.html
var views embed.FS

//go:embed views/layouts/*.html
var layouts embed.FS

//go:embed views/partials/*.html
var partials embed.FS

type TemplateConfig struct {
	Dev         bool
	CustomFuncs template.FuncMap
}

type DevTemplate struct {
	funcs template.FuncMap
}

func (d *DevTemplate) ExecuteTemplate(w io.Writer, name string, data any) error {
	t := template.New("default.html").Funcs(d.funcs)
	t, err := t.ParseFS(os.DirFS("internal/ui/views/layouts"), "*.html")
	if err != nil {
		return err
	}

	t, err = t.ParseFS(os.DirFS("internal/ui/views/partials"), "*.html")
	if err != nil {
		return err
	}

	p := parseTemplateName(name)

	t, err = t.ParseFiles("internal/ui/views/" + p.File)
	if err != nil {
		return err
	}

	return t.ExecuteTemplate(w, p.Template, data)
}

type EmbeddedTemplate struct {
	shared *template.Template
}

func (e *EmbeddedTemplate) ExecuteTemplate(w io.Writer, name string, data any) error {
	shared, err := e.shared.Clone()
	if err != nil {
		return err
	}

	p := parseTemplateName(name)

	t, err := shared.ParseFS(views, "views/"+p.File)
	if err != nil {
		return err
	}

	return t.ExecuteTemplate(w, p.Template, data)
}

func NewTemplate(conf TemplateConfig) *echo.TemplateRenderer {
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
		shared = template.Must(shared.ParseFS(layouts, "views/layouts/*.html"))
		shared = template.Must(shared.ParseFS(partials, "views/partials/*.html"))
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

type templateNameParts struct {
	File     string
	Template string
}

func parseTemplateName(name string) templateNameParts {
	p := templateNameParts{
		File:     name,
		Template: name,
	}

	if strings.Contains(name, "#") {
		parts := strings.SplitN(name, "#", 2)
		p.File = parts[0]
		p.Template = parts[1]
	}

	return p
}
