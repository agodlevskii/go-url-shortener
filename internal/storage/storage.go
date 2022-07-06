package storage

type ShortURL struct {
	ID  string
	URL string
	UID string
}

type Storager interface {
	Add(batch []ShortURL) ([]ShortURL, error)
	Has(id string) (bool, error)
	Get(id string) (string, error)
	GetAll(userID string) ([]ShortURL, error)
	Clear()
	Ping() bool
}
