package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type errorViewData struct {
	View string
	Code int
	Err  error
}

func newErrorViewData(err error) errorViewData {
	var h *echo.HTTPError
	if !errors.As(err, &h) {
		h = echo.NewHTTPError(http.StatusInternalServerError).WithInternal(err)
	}

	return errorViewData{
		View: "error",
		Code: h.Code,
		Err:  h,
	}
}

func customNotFoundHandler(c echo.Context) error {
	data := newErrorViewData(
		echo.NewHTTPError(http.StatusNotFound, "Page not found"),
	)

	return c.Render(http.StatusNotFound, "404.html", data)
}

func apiNotFoundHandler(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)

	if c.Response().Committed {
		return
	}

	if strings.HasPrefix(c.Request().URL.Path, "/api/") {
		c.Echo().DefaultHTTPErrorHandler(err, c)
		return
	}

	data := newErrorViewData(err)

	if data.Code == http.StatusNotFound {
		err = c.Render(http.StatusNotFound, "404.html", data)
	} else {
		err = c.Render(data.Code, "error.html", data)
	}

	if err != nil {
		c.Logger().Error(fmt.Errorf("while rendering error page, encountered another error: %w", err))
	}
}
