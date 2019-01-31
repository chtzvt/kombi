package db

type ShortlinkReader interface {
	GetLinkDestination(id string) (dest string, err error)
}

type CachedShortlinkReader interface {
	ShortlinkReader
	GetLinkDestinationCached(id string) (dest string, wasCached bool, err error)
	GetLinkDestinationUncached(id string) (dest string, err error)
}

type ShortlinkWriter interface {
	CreateLink(id string, dest string) (err error)
	UpdateLinkDestination(id string, dest string) (err error)
	DeleteLink(id string) (err error)
}

type CachedShortlinkWriter interface {
	ShortlinkWriter
	AddLinkToCache(id string, dest string) (err error)
	DeleteLinkFromCache(id string) (err error)
	DeleteExpiredEntries() (err error)
	FlushCache() (err error)
}

type ShortlinkReadWriter interface {
	ShortlinkReader
	ShortlinkWriter
}

