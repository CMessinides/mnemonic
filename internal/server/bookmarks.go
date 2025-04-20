package server

import (
	"net/http"

	"github.com/cmessinides/mnemonic/internal/bookmark"
	"github.com/labstack/echo/v4"
)

type bookmarksAPI struct {
	store bookmark.BookmarkStore
}

func (a *bookmarksAPI) Create(c echo.Context) error {
	var init struct {
		Title string
		URL   string
		Tags  []string
	}

	init.Tags = make([]string, 0)

	err := echo.FormFieldBinder(c).
		MustString("title", &init.Title).
		MustString("url", &init.URL).
		Strings("tags", &init.Tags).
		BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	b, err := a.store.Create(init.Title, init.URL, init.Tags)
	if err != nil {
		if bookmark.IsURLExists(err) {
			return echo.NewHTTPError(http.StatusBadRequest, "URL already exists")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(http.StatusCreated, b)
}

func (a *bookmarksAPI) List(c echo.Context) error {
	var opts struct {
		Limit  uint
		Offset uint
	}

	err := echo.QueryParamsBinder(c).
		Uint("limit", &opts.Limit).
		Uint("offset", &opts.Offset).
		BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if opts.Limit == 0 {
		opts.Limit = 100
	}

	bookmarks, err := a.store.GetMany(opts.Limit, opts.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, bookmarks)
}
