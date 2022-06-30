package storage

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestDBRepo_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := DBRepo{db: db}
	tests := getAddTestCases(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO urls").WithArgs(tt.args.id, tt.args.url, UserID)

			tt.repo.Add(UserID, tt.args.id, tt.args.url)
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}

	repo.Clear()
}

func TestDBRepo_AddAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := DBRepo{db: db}
	tests := getAddAllTestCases(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectPrepare("INSERT INTO urls")
			for id, url := range tt.batch {
				mock.ExpectExec("INSERT INTO urls").WithArgs(id, url, UserID).WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mock.ExpectCommit()

			tt.repo.AddAll(UserID, tt.batch)
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}

	repo.Clear()
}

func TestDBRepo_Clear(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := DBRepo{db: db}
	tests := getClearTestCases(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("DELETE FROM urls")
			tt.repo.Clear()
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func getTestDBRepo() (DBRepo, error) {
	db, _, err := sqlmock.New()
	if err != nil {
		return DBRepo{}, err
	}
	return DBRepo{db: db}, nil
}
