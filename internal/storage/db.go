package storage

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
)

type DBRepo struct {
	db *sql.DB
}

func NewDBRepo(url string) (DBRepo, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return DBRepo{}, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (id VARCHAR(10), url VARCHAR(255), uid VARCHAR(16), UNIQUE(id), UNIQUE (url))")
	if err != nil {
		return DBRepo{}, err
	}
	return DBRepo{db: db}, nil
}

func (repo DBRepo) Add(batch []ShortURL) ([]ShortURL, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(`INSERT INTO urls(id, url, uid) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res := make([]ShortURL, len(batch))
	for i, sURL := range batch {
		var newID string

		err = stmt.QueryRow(sURL.ID, sURL.URL, sURL.UID).Scan(&newID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) || err.Error() == "sql: no rows in result set" {
				err = repo.db.QueryRow("SELECT id FROM urls WHERE url = $1", sURL.URL).Scan(&newID)
			}

			if err != nil {
				log.Error(err)
				if err = tx.Rollback(); err != nil {
					log.Error("unable to rollback: ", err)
				}

				return nil, err
			}
		}

		res[i] = ShortURL{
			ID:  newID,
			URL: sURL.URL,
			UID: sURL.UID,
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error("unable to commit: ", err)
		return nil, err
	}

	return res, nil
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

func (repo DBRepo) GetAll(userID string) ([]ShortURL, error) {
	rows, err := repo.db.Query("SELECT * FROM urls WHERE uid = $1", userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	urls := make([]ShortURL, 0)
	for rows.Next() {
		var sURL ShortURL
		err = rows.Scan(&sURL.ID, &sURL.URL, &sURL.UID)
		if err != nil {
			return nil, err
		}

		urls = append(urls, sURL)
	}

	return urls, nil
}

func (repo DBRepo) Clear() {
	repo.db.Exec("DELETE FROM urls")
}

func (repo DBRepo) Ping() bool {
	return repo.db.Ping() == nil
}
