package storage

type ShortURL struct {
	ID      string
	URL     string
	UID     string
	Deleted bool
}

type Storager interface {
	Add(batch []ShortURL) ([]ShortURL, error)
	Clear()
	Delete(batch []ShortURL) error
	Get(id string) (ShortURL, error)
	GetAll(userID string) ([]ShortURL, error)
	Has(id string) (bool, error)
	Ping() bool
}
