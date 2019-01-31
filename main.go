package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"kombi/db"
	"log"
	"net/http"
	"time"
)

var linkDb = &db.CachedSqliteDb {
	FilePath:               "./links.db",
	CacheDefaultExpiration: 5*time.Minute,
	CacheExpiredPurgeTime:  5*time.Minute,
	ShouldInit:             false,
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())


	if err := db.CachedSqliteDatabase(linkDb); err != nil {
		log.Fatal("Error creating link database: ", err.Error())
	}

	// Routes
	e.GET("/", LinkServer)
	//e.GET("/settings", Settings)
	//e.GET("/api", API)

	e.GET("/:link", LinkServer)

	time.Sleep(2*time.Second)

	if err := linkDb.CreateLink("", "https://www.hackucf.org/"); err != nil {
		log.Printf("Error adding link to database: %s", err.Error())
	}

	if err := linkDb.CreateLink("ct2", "http://chill.ctis.me/"); err != nil {
		log.Printf("Error adding link to database: %s", err.Error())
	}

	if err := linkDb.CreateLink("ct", "https://www.ctis.me/"); err != nil {
		log.Printf("Error adding link to database: %s", err.Error())
	}

	if err := linkDb.DeleteLink("ct2"); err != nil {
		log.Printf("Error removing link from database: %s", err.Error())
	}

	// Start server
	e.Logger.Fatal(e.Start(":80"))
}


func LinkServer(c echo.Context) error {
	if dest, err := linkDb.GetLinkDestination(c.Param("link")); err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("\"%s\": not found: %s (%s)", c.Param("link"), dest, err.Error()))
	} else {
		return c.Redirect(http.StatusFound, dest)
	}
}

