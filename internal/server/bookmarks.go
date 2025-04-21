package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

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

	init.Tags = []string{}

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
	var page, pageSize uint64
	err := echo.QueryParamsBinder(c).
		Uint64("page", &page).
		Uint64("pageSize", &pageSize).
		BindError()
	if err != nil {
		var berr *echo.BindingError
		if errors.As(err, &berr) {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Sprintf(
					"invalid value for %s (got: %s)", berr.Field, strings.Join(berr.Values, ", "),
				),
			)
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 25
	} else if pageSize > 100 {
		pageSize = 100
	}

	bp, err := a.store.GetPage(page, pageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, bp)
}
