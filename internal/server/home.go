package server

import (
	"net/http"

	"github.com/cmessinides/mnemonic/internal/bookmark"
	"github.com/cmessinides/mnemonic/internal/pagination"
	"github.com/labstack/echo/v4"
)

type homeController struct {
	bookmarks bookmark.BookmarkStore
}

func (h *homeController) Show(c echo.Context) error {
	status := http.StatusOK
	var data struct {
		View           string
		Bookmarks      *pagination.Page[*bookmark.Bookmark]
		BookmarksError string
	}
	data.View = "home"

	bookmarks, err := h.bookmarks.GetPage(1, 30)
	if err != nil {
		c.Logger().Warn(err)
		status = http.StatusInternalServerError
		data.BookmarksError = err.Error()
	} else {
		data.Bookmarks = bookmarks
	}

	return c.Render(status, "home.html", data)
}
