package storage

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDBRepo_Add(t *testing.T) {
	for _, tt := range getAddTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer func(db *sql.DB) {
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}(db)

			r := DBRepo{db: db}
			coverInitExpect(mock, tt.state)
			mock.ExpectClose()
			got, err := r.Add(context.Background(), tt.state)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want.id, got[0].ID)
		})
	}
}

func TestDBRepo_Clear(t *testing.T) {
	for _, tt := range getClearTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer func(db *sql.DB) {
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}(db)

			r := DBRepo{db: db}
			q := regexp.QuoteMeta(AddURLs)
			mock.ExpectBegin()
			mock.ExpectPrepare(q)
			for _, v := range tt.state {
				mock.ExpectQuery(q).
					WithArgs(v.ID, v.URL, v.UID, v.Deleted).
					WillReturnRows(mock.NewRows([]string{"id"}).AddRow(v.ID))
			}
			mock.ExpectCommit()
			mock.ExpectExec(DeleteURL).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectClose()

			if _, err := r.Add(context.Background(), tt.state); err != nil {
				t.Fatal(err)
			}
			r.Clear(context.Background())
		})
	}
}

func TestDBRepo_Delete(t *testing.T) {
	for _, tt := range getDeleteTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer func(db *sql.DB) {
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}(db)

			r := DBRepo{db: db}
			ids := make([]string, len(tt.state))
			for i, sURL := range tt.state {
				ids[i] = sURL.ID
			}

			coverInitExpect(mock, tt.state)
			mock.ExpectExec(regexp.QuoteMeta(DeleteUserURLs)).
				WithArgs(UserID, pq.Array(ids)).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectClose()

			if _, err := r.Add(context.Background(), tt.state); err != nil {
				t.Fatal(err)
			}

			err := r.Delete(context.Background(), tt.state)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDBRepo_Get(t *testing.T) {
	for _, tt := range getGetTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer func(db *sql.DB) {
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}(db)

			r := DBRepo{db: db}
			var res ShortURL
			for _, v := range tt.state {
				if v.ID == tt.id {
					res = v
				}
			}

			coverInitExpect(mock, tt.state)
			eq := mock.ExpectQuery(regexp.QuoteMeta(GetURL)).WithArgs(tt.id)
			if tt.want != "" {
				rows := sqlmock.NewRows([]string{"id", "url", "uid", "deleted"}).
					AddRow(res.ID, res.URL, res.UID, res.Deleted)
				eq.WillReturnRows(rows)
			} else {
				eq.WillReturnError(sql.ErrNoRows)
			}
			mock.ExpectClose()

			if _, err := r.Add(context.Background(), tt.state); err != nil {
				t.Fatal(err)
			}
			got, err := r.Get(context.Background(), tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got.URL)
		})
	}
}

func TestDBRepo_GetAll(t *testing.T) {
	for _, tt := range getGetAllTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer func(db *sql.DB) {
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}(db)

			r := DBRepo{db: db}
			ids := make([]string, len(tt.state))
			for i, sURL := range tt.state {
				ids[i] = sURL.ID
			}

			rows := sqlmock.NewRows([]string{"id", "url", "uid", "deleted"})
			for _, v := range tt.state {
				if tt.want[v.ID] {
					rows.AddRow(v.ID, v.URL, v.UID, v.Deleted)
				}
			}

			coverInitExpect(mock, tt.state)
			mock.ExpectQuery(regexp.QuoteMeta(GetUserURLs)).WithArgs(UserID).WillReturnRows(rows)
			mock.ExpectClose()

			if _, err := r.Add(context.Background(), tt.state); err != nil {
				t.Fatal(err)
			}
			got, err := r.GetAll(context.Background(), UserID)
			gotMap := make(map[string]bool)
			for _, gv := range got {
				gotMap[gv.ID] = true
			}

			assert.Equal(t, tt.wantErr, err != nil)
			for id, want := range tt.want {
				assert.Equal(t, want, gotMap[id])
			}
			for _, tv := range got {
				assert.True(t, tt.want[tv.ID])
			}
		})
	}
}

func TestDBRepo_Has(t *testing.T) {
	for _, tt := range getHasTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer func(db *sql.DB) {
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}(db)
			r := DBRepo{db: db}
			var exp int64
			if tt.want {
				exp = 1
			}

			coverInitExpect(mock, tt.state)
			mock.ExpectQuery(regexp.QuoteMeta(HasURL)).
				WithArgs(tt.id).
				WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(exp))
			mock.ExpectClose()

			if _, err := r.Add(context.Background(), tt.state); err != nil {
				t.Fatal(err)
			}
			got, err := r.Has(context.Background(), tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func getMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	return db, mock
}

func coverInitExpect(mock sqlmock.Sqlmock, state []ShortURL) {
	q := regexp.QuoteMeta(AddURLs)
	mock.ExpectBegin()
	mock.ExpectPrepare(q)
	for _, v := range state {
		mock.ExpectQuery(q).
			WithArgs(v.ID, v.URL, v.UID, v.Deleted).
			WillReturnRows(mock.NewRows([]string{"id"}).AddRow(v.ID))
	}
	mock.ExpectCommit()
}
