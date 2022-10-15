package storage

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v4/stdlib" // SQL driver
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// DBRepo describes the SQL implementation of the Storager interface.
type DBRepo struct {
	db *sql.DB
}

// NewDBRepo returns a new instance of the DBRepo type.
// If the DB didn't connect, or the DB table creation has failed, the error will be returned.
func NewDBRepo(url string) (DBRepo, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return DBRepo{}, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (id VARCHAR(10), url VARCHAR(255), uid VARCHAR(16), deleted boolean, UNIQUE(id), UNIQUE(url))")
	if err != nil {
		return DBRepo{}, err
	}
	return DBRepo{db: db}, nil
}

// Add provides a functionality to save a slice of the ShortURL data into the SQL repository.
// If the insert fails, or the saved data fails to return, the changes will be rollback, and the error will be returned.
func (repo DBRepo) Add(batch []ShortURL) ([]ShortURL, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(`INSERT INTO urls(id, url, uid, deleted) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res := make([]ShortURL, len(batch))
	for i, sURL := range batch {
		var newID string

		err = stmt.QueryRow(sURL.ID, sURL.URL, sURL.UID, sURL.Deleted).Scan(&newID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
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

// Has checks if the repository contains the ShortURL with a specific ID.
// If the select query fails, the error will be returned.
func (repo DBRepo) Has(id string) (bool, error) {
	var cnt int64
	err := repo.db.QueryRow("SELECT COUNT(*) FROM urls WHERE id = $1", id).Scan(&cnt)
	return cnt != 0, err
}

// Get returns the ShortURL value by its ID.
// If the select query fails, the error will be returned.
func (repo DBRepo) Get(id string) (ShortURL, error) {
	var sURL ShortURL
	err := repo.db.QueryRow("SELECT * FROM urls WHERE id = $1", id).Scan(&sURL.ID, &sURL.URL, &sURL.UID, &sURL.Deleted)
	return sURL, err
}

// GetAll returns all the ShortURL values created by the specified user.
// If the repository doesn't have any associated value, the empty slice will be returned.
// If the select query fails, the error will be returned.
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
		err = rows.Scan(&sURL.ID, &sURL.URL, &sURL.UID, &sURL.Deleted)
		if err != nil {
			return nil, err
		}

		urls = append(urls, sURL)
	}

	return urls, nil
}

// Clear marks all existing values in the repository as deleted.
func (repo DBRepo) Clear() {
	if _, err := repo.db.Exec("UPDATE urls SET deleted = true"); err != nil {
		log.Error(err)
	}
}

// Ping provides a proxy for the sql package's ping functionality.
func (repo DBRepo) Ping() bool {
	return repo.db.Ping() == nil
}

// Delete marks all specified ShortURL values in repository as deleted.
// The deletion of the value is available only for its owner. All other values will be skipped.
// If the update query fails, the error will be returned.
func (repo DBRepo) Delete(batch []ShortURL) error {
	if len(batch) == 0 {
		return nil
	}

	userID := batch[0].UID
	ids := make([]string, len(batch))
	for i, sURL := range batch {
		ids[i] = sURL.ID
	}

	_, err := repo.db.Exec("UPDATE urls SET deleted = true WHERE uid = $1 AND id = any($2)", userID, pq.Array(ids))
	return err
}

func (repo DBRepo) Close() error {
	return repo.db.Close()
}
