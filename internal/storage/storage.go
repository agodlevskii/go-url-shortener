package storage

type URLRes struct {
	url string
	uid string
}

type Storager interface {
	Add(userID, id, url string) error
	AddAll(userId string, batch map[string]string) error
	Has(id string) (bool, error)
	Get(id string) (string, error)
	GetAll(userID string) (map[string]string, error)
	Clear()
}
