package linkserver

import (
	"github.com/labstack/echo"
	"kombi/db"
	"net/http"
)

var linkDb db.ShortlinkReader

func RegisterDatabase(database db.ShortlinkReader) {
	linkDb = database
}

func RegisterPathHandlers(e *echo.Echo) {
	e.GET("/", LinkServer)
	e.GET("/:link", LinkServer)
}

func LinkServer(c echo.Context) error {
	if dest, err := linkDb.GetLinkDestination(c.Param("link")); err != nil {
		return c.Redirect(http.StatusFound, "/")
	} else {
		return c.Redirect(http.StatusFound, dest)
	}
}
