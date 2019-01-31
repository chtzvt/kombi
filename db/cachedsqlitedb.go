package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	goCache "github.com/patrickmn/go-cache"
	"time"
)

type CachedSqliteDb struct {
	FilePath               string
	CacheDefaultExpiration time.Duration
	CacheExpiredPurgeTime  time.Duration // AKA Cleanup interval- see go-cache janitor docs
	Initialize             bool
	cache                  *goCache.Cache
	sqlite                 *sql.DB
}

func CachedSqliteDatabase(db *CachedSqliteDb) error {
	sqliteDb, err := sql.Open("sqlite3", db.FilePath)
	if err != nil {
		return errors.New("CachedSqliteDatabase: " + err.Error())
	}

	db.sqlite = sqliteDb
	db.cache = goCache.New(db.CacheDefaultExpiration, db.CacheExpiredPurgeTime)
	if err = db.initialize(); err != nil {
		return errors.New("CachedSqliteDatabase: " + err.Error())
	}

	return nil
}

func (db *CachedSqliteDb) initialize() (err error) {
	if db.Initialize == false {
		return nil
	}

	tableSetup := `
	CREATE TABLE shortlinks (shortpath TEXT UNIQUE, destination TEXT, hits INT, created TEXT);
	CREATE INDEX shortlink_shortpaths ON shortlinks (shortpath);
	CREATE INDEX shortlink_destinations ON shortlinks (destination);
	CREATE INDEX shortlink_shortpath_dest ON shortlinks (shortpath, destination);
	`

	_, err = db.sqlite.Exec(tableSetup)
	if err != nil {
		return errors.New("initialize: " + err.Error())
	}

	return nil
}

func (db *CachedSqliteDb) GetLinkDestination(id string) (dest string, err error) {
	dest, wasCached, err := db.GetLinkDestinationCached(id)
	if wasCached {
		fmt.Printf("HIT for link id [%s], dest: %s ", id, dest)
		return dest, nil
	} else {
		fmt.Printf("MISS for link id [%s] ", id)
		dest, err = db.GetLinkDestinationUncached(id)
		if err != nil && err != sql.ErrNoRows {
			return dest, errors.New("GetLinkDestination: " + err.Error())
		}

		_ = db.AddLinkToCache(id, dest)

		return
	}
}

func (db *CachedSqliteDb) GetLinkDestinationUncached(id string) (dest string, err error) {
	row := db.sqlite.QueryRow(`SELECT destination FROM shortlinks WHERE shortpath = ?;`, id)

	if err = row.Scan(&dest); err != nil {
		dest = "/"
	}

	return
}

func (db *CachedSqliteDb) CreateLink(id string, dest string) (err error) {
	result, err := db.sqlite.Exec(`INSERT INTO shortlinks(shortpath, destination, hits, created) VALUES(?, ?, ?, ?)`, id, dest, 0, time.Now().Unix())
	if err != nil {
		return errors.New("CreateLink: " + err.Error())
	}

	if nRows, _ := result.RowsAffected(); nRows < 1 {
		return errors.New("CreateLink: inserted records into database, but no rows were affected.")
	}

	_ = db.AddLinkToCache(id, dest)

	return nil
}

func (db *CachedSqliteDb) UpdateLinkDestination(id string, dest string) (err error) {
	result, err := db.sqlite.Exec(`UPDATE shortlinks SET destination = ? WHERE shortpath = ?`, dest, id)
	if err != nil {
		return errors.New("UpdateLinkDestination: " + err.Error())
	}

	if nRows, _ := result.RowsAffected(); nRows < 1 {
		return errors.New("UpdateLinkDestination: updated record in database, but no rows were affected.")
	}

	_ = db.AddLinkToCache(id, dest)

	return
}

func (db *CachedSqliteDb) DeleteLink(id string) (err error) {

	result, err := db.sqlite.Exec(`DELETE FROM shortlinks WHERE shortpath = ?`, id)
	if err != nil {
		return errors.New("DeleteLink: " + err.Error())
	}

	if nRows, _ := result.RowsAffected(); nRows < 1 {
		return errors.New("DeleteLink: deleted records from database, but no rows were affected.")
	}

	_ = db.DeleteLinkFromCache(id)

	return
}

func (db *CachedSqliteDb) AddLinkToCache(id string, dest string) (err error) {
	db.cache.Set(id, dest, 0)
	return
}

func (db *CachedSqliteDb) GetLinkDestinationCached(id string) (dest string, wasCached bool, err error) {
	destination, wasCached := db.cache.Get(id) // Never errors
	if wasCached == false {
		// Type assertion will fail if destination is nil, so we'll need to
		// fill it in with an empty value if that's the case.
		destination = ""
	}

	return destination.(string), wasCached, nil
}

func (db *CachedSqliteDb) DeleteLinkFromCache(id string) (err error) {
	db.cache.Delete(id) // Never errors
	return
}

func (db *CachedSqliteDb) DeleteExpiredEntries() (err error) {
	db.cache.DeleteExpired() // Never errors
	return
}

func (db *CachedSqliteDb) FlushCache() (err error) {
	db.cache.Flush() // Never errors
	return
}
