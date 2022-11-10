package storage

import (
	"context"
	"database/sql"
	"errors"
	"go-url-shortener/internal/apperrors"

	_ "github.com/jackc/pgx/v4/stdlib" // SQL driver
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	CreateURLTable = `CREATE TABLE IF NOT EXISTS urls(
    	id VARCHAR(10),
    	url VARCHAR(255),
    	uid VARCHAR(16),
    	deleted boolean,
    	UNIQUE(id), UNIQUE(url))`
	AddURLs = `INSERT INTO urls(id, url, uid, deleted) VALUES ($1, $2, $3, $4)
                                        ON CONFLICT DO NOTHING RETURNING id`
	HasURL         = `SELECT COUNT(*) FROM urls WHERE id = $1`
	GetURLID       = `SELECT id FROM urls WHERE url = $1`
	GetURL         = `SELECT * FROM urls WHERE id = $1`
	GetUserURLs    = `SELECT * FROM urls WHERE uid = $1`
	DeleteURL      = `UPDATE urls u SET deleted = true WHERE u.id <> '' IS NOT TRUE`
	DeleteUserURLs = `UPDATE urls SET deleted = true WHERE uid = $1 AND id = any($2)`
)

// DBRepo describes the SQL implementation of the Storager interface.
type DBRepo struct {
	db *sql.DB
}

// NewDBRepo returns a new instance of the DBRepo type.
// If the DB didn't connect, or the DB table creation has failed, the error will be returned.
func NewDBRepo(ctx context.Context, url string) (DBRepo, error) {
	if url == "" {
		return DBRepo{}, errors.New(apperrors.EmptyDBURL)
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return DBRepo{}, err
	}

	_, err = db.ExecContext(ctx, CreateURLTable)
	if err != nil {
		return DBRepo{}, err
	}
	return DBRepo{db: db}, nil
}

// Add provides a functionality to save a slice of the ShortURL data into the SQL repository.
// If the insert fails, or the saved data fails to return, the changes will be rollback, and the error will be returned.
func (repo DBRepo) Add(ctx context.Context, batch []ShortURL) ([]ShortURL, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, AddURLs)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		if cErr := stmt.Close(); cErr != nil {
			log.Error(cErr)
		}
	}(stmt)

	res := make([]ShortURL, len(batch))
	for i, sURL := range batch {
		var newID string

		err = stmt.QueryRow(sURL.ID, sURL.URL, sURL.UID, sURL.Deleted).Scan(&newID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = repo.db.QueryRowContext(ctx, GetURLID, sURL.URL).Scan(&newID)
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
func (repo DBRepo) Has(ctx context.Context, id string) (bool, error) {
	var cnt int64
	err := repo.db.QueryRowContext(ctx, HasURL, id).Scan(&cnt)
	return cnt != 0, err
}

// Get returns the ShortURL value by its ID.
// If the select query fails, the error will be returned.
func (repo DBRepo) Get(ctx context.Context, id string) (ShortURL, error) {
	var sURL ShortURL
	err := repo.db.QueryRowContext(ctx, GetURL, id).Scan(&sURL.ID, &sURL.URL, &sURL.UID, &sURL.Deleted)
	return sURL, err
}

// GetAll returns all the ShortURL values created by the specified user.
// If the repository doesn't have any associated value, the empty slice will be returned.
// If the select query fails, the error will be returned.
func (repo DBRepo) GetAll(ctx context.Context, userID string) ([]ShortURL, error) {
	rows, err := repo.db.QueryContext(ctx, GetUserURLs, userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer func(rows *sql.Rows) {
		if cErr := rows.Close(); cErr != nil {
			log.Error(cErr)
		}
	}(rows)

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
func (repo DBRepo) Clear(ctx context.Context) {
	if _, err := repo.db.ExecContext(ctx, DeleteURL); err != nil {
		log.Error(err)
	}
}

// Ping provides a proxy for the sql package's ping functionality.
func (repo DBRepo) Ping(ctx context.Context) bool {
	return repo.db.PingContext(ctx) == nil
}

// Delete marks all specified ShortURL values in repository as deleted.
// The deletion of the value is available only for its owner. All other values will be skipped.
// If the update query fails, the error will be returned.
func (repo DBRepo) Delete(ctx context.Context, batch []ShortURL) error {
	if len(batch) == 0 {
		return nil
	}

	userID := batch[0].UID
	ids := make([]string, 0, len(batch))
	for _, sURL := range batch {
		ids = append(ids, sURL.ID)
	}

	_, err := repo.db.ExecContext(ctx, DeleteUserURLs, userID, pq.Array(ids))
	return err
}

func (repo DBRepo) Close() error {
	return repo.db.Close()
}
