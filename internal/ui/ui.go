package ui

import (
	"html/template"
	"io/fs"
	"maps"

	"github.com/labstack/echo/v4"
)

type UI struct {
	Assets   fs.FS
	Renderer echo.Renderer
	config   UIConfig
}

type UIConfig struct {
	Dev           bool
	AssetPath     string
	TemplateFuncs template.FuncMap
}

func (c UIConfig) WithTemplateFuncs(funcs template.FuncMap) UIConfig {
	var merged template.FuncMap
	if c.TemplateFuncs == nil {
		merged = template.FuncMap{}
	} else {
		merged = maps.Clone(c.TemplateFuncs)
	}

	maps.Copy(merged, funcs)

	return UIConfig{
		Dev:           c.Dev,
		AssetPath:     c.AssetPath,
		TemplateFuncs: merged,
	}
}

func NewUI(conf UIConfig) *UI {
	a := NewAssetsFS(conf)

	conf = conf.WithTemplateFuncs(a.TemplateFuncs())
	t := NewTemplateRenderer(conf)

	return &UI{
		Assets:   a,
		Renderer: t,
		config:   conf,
	}
}

func (u *UI) ConfigureServer(e *echo.Echo) {
	e.StaticFS(u.config.AssetPath, u.Assets)
	e.Renderer = u.Renderer
}
