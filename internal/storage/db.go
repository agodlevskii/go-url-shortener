package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBRepo struct {
	db *sql.DB
}

type DBURLRes struct {
	id  string
	url string
	uid string
}

func NewDBRepo(url string) (DBRepo, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return DBRepo{}, err
	}

	return DBRepo{db: db}, nil
}

func (repo DBRepo) Add(userID, id, url string) error {
	res, err := repo.db.Exec("INSERT INTO urls(id, url, uid) VALUES ($1, $2, $3)", id, url, userID)
	if err != nil || res == nil {
		return err
	}
	return nil
}

func (repo DBRepo) Has(id string) (bool, error) {
	var cnt int64
	err := repo.db.QueryRow("SELECT COUNT(*) FROM urls WHERE id = $1", id).Scan(&cnt)
	return cnt != 0, err
}

func (repo DBRepo) Get(id string) (string, error) {
	var url string
	err := repo.db.QueryRow("SELECT url FROM urls WHERE id = $1", id).Scan(&url)
	return url, err
}

func (repo DBRepo) GetAll(userID string) (map[string]string, error) {
	rows, err := repo.db.Query("SELECT * FROM urls WHERE uid = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make(map[string]string)
	for rows.Next() {
		var res DBURLRes
		err = rows.Scan(&res.id, &res.url, &res.uid)
		if err != nil {
			return nil, err
		}

		urls[res.id] = res.url
	}

	return urls, nil
}

func (repo DBRepo) Clear() {
	repo.db.QueryRow("DELETE FROM urls")
}
