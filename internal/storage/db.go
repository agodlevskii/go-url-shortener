package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (id VARCHAR(10), url VARCHAR(255), uid VARCHAR(16), UNIQUE(id))")
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

func (repo DBRepo) AddAll(userID string, batch map[string]string) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO urls(id, url, uid) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id, url := range batch {
		if _, err = stmt.Exec(id, url, userID); err != nil {
			log.Error(err)

			if err = tx.Rollback(); err != nil {
				log.Fatal("unable to rollback: ", err)
			}

			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Fatal("unable to commit: ", err)
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
	if rows.Err() != nil {
		return nil, rows.Err()
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
	repo.db.Exec("DELETE FROM urls")
}
