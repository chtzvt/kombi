package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"kombi/db"
	"kombi/linkserver"
	"log"
	"time"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	linkDb := &db.CachedSqliteDb{
		FilePath:               "./links.db",
		CacheDefaultExpiration: 5 * time.Minute,
		CacheExpiredPurgeTime:  5 * time.Minute,
		Initialize:             false,
	}

	if err := db.CachedSqliteDatabase(linkDb); err != nil {
		log.Fatal(fmt.Sprintf("Couldn't configure database: ", err.Error()))
	}

	linkserver.RegisterDatabase(linkDb)
	linkserver.RegisterPathHandlers(e)

	//e.GET("/kombi", Settings)
	//e.POST("/kombi/api", API)

	/*
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
	*/

	e.Logger.Fatal(e.Start(":8080"))
}
