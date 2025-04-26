package server

import (
	"errors"
	"fmt"
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

	init.Tags = []string{}

	err := echo.FormFieldBinder(c).
		MustString("title", &init.Title).
		MustString("url", &init.URL).
		Strings("tags", &init.Tags).
		BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).WithInternal(err)
	}

	b, err := a.store.Create(init.Title, init.URL, init.Tags)
	if err != nil {
		return fail(err)
	}

	return c.JSON(http.StatusCreated, b)
}

func (a *bookmarksAPI) Read(c echo.Context) error {
	var id int64

	err := echo.PathParamsBinder(c).
		MustInt64("id", &id).
		BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required").WithInternal(err)
	}

	b, err := a.store.Get(id)
	if err != nil {
		return fail(err)
	}

	return c.JSON(http.StatusOK, b)
}

func (a *bookmarksAPI) Update(c echo.Context) error {
	var id int64
	var title string
	var url string
	var archived bool
	var tags []string

	err := echo.PathParamsBinder(c).
		MustInt64("id", &id).
		BindError()
	if err != nil {
		return err
	}

	err = echo.FormFieldBinder(c).
		String("title", &title).
		String("url", &url).
		Bool("archived", &archived).
		Strings("tags", &tags).
		BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).WithInternal(err)
	}

	err = a.store.Update(bookmark.BookmarkPatch{
		ID:       id,
		Title:    &title,
		URL:      &url,
		Archived: &archived,
		Tags:     tags,
	})
	if err != nil {
		return fail(err)
	}

	return c.NoContent(http.StatusOK)
}

func (a *bookmarksAPI) List(c echo.Context) error {
	var page, pageSize uint64
	err := echo.QueryParamsBinder(c).
		Uint64("page", &page).
		Uint64("pageSize", &pageSize).
		BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).WithInternal(err)
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
		return fail(err)
	}

	return c.JSON(http.StatusOK, bp)
}

func (a *bookmarksAPI) Delete(c echo.Context) error {
	var id int64

	err := echo.PathParamsBinder(c).
		MustInt64("id", &id).
		BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required").WithInternal(err)
	}

	err = a.store.Delete(id)
	if err != nil {
		return fail(err)
	}

	return c.NoContent(http.StatusOK)
}

func fail(err error) *echo.HTTPError {
	if bookmark.IsNotFound(err) {
		return echo.NewHTTPError(http.StatusNotFound, "bookmark not found").WithInternal(err)
	}

	var ue *bookmark.URLExistsError
	if errors.As(err, &ue) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("a bookmark already exists with that URL (%s)", ue.URL)).WithInternal(err)
	}

	return echo.NewHTTPError(http.StatusInternalServerError).WithInternal(err)
}
