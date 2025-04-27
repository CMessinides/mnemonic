package server

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"

	"github.com/labstack/echo/v4"
)

//go:embed all:public/assets
var assets embed.FS

type AssetConfig struct {
	PublicPath string
	Dev        bool
}

type AssetsFS struct {
	PublicPath string
	fs         fs.FS
}

func NewAssetsFS(conf AssetConfig) *AssetsFS {
	a := &AssetsFS{
		PublicPath: conf.PublicPath,
	}

	if conf.Dev {
		a.fs = os.DirFS("internal/server/public/assets")
	} else {
		a.fs = echo.MustSubFS(assets, "public/assets")
	}

	return a
}

func (a *AssetsFS) Open(name string) (fs.File, error) {
	return a.fs.Open(name)
}

func (a *AssetsFS) FileExists(name string) (bool, error) {
	sf, ok := a.fs.(fs.StatFS)
	if !ok {
		return false, errors.New("assets filesystem is not a StatFS")
	}

	info, err := sf.Stat(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}

	return !info.IsDir(), nil
}

func (a *AssetsFS) TemplateFuncs() template.FuncMap {
	asset := func(filename string) string {
		return a.PublicPath + "/" + filename
	}

	stylesheet := func(filename string) template.HTML {
		return template.HTML(
			fmt.Sprintf(
				`<link rel="stylesheet" href="%s/%s">`,
				template.HTMLEscapeString(a.PublicPath),
				template.HTMLEscapeString(filename),
			),
		)
	}

	script := func(filename string) template.HTML {
		return template.HTML(
			fmt.Sprintf(
				`<script defer type="module" src="%s/%s"></script>`,
				template.HTMLEscapeString(a.PublicPath),
				template.HTMLEscapeString(filename),
			),
		)
	}

	return template.FuncMap{
		"asset":      asset,
		"stylesheet": stylesheet,
		"script":     script,
		"assetIfExists": func(filename string) string {
			if exists, _ := a.FileExists(filename); !exists {
				return ""
			}

			return asset(filename)
		},
		"stylesheetIfExists": func(filename string) template.HTML {
			if exists, _ := a.FileExists(filename); !exists {
				return ""
			}

			return stylesheet(filename)
		},
		"scriptIfExists": func(filename string) template.HTML {
			if exists, _ := a.FileExists(filename); !exists {
				return ""
			}

			return script(filename)
		},
		"icon": func(id string) template.HTML {
			return template.HTML(
				fmt.Sprintf(
					`<svg class="icon" height="16" width="16"><use xlink:href="%s/icons.svg#%s"></use></svg>`,
					a.PublicPath,
					id,
				),
			)
		},
	}
}
