package storage

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestDBRepo_Add(t *testing.T) {
	for _, tt := range getAddTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer db.Close()
			r := DBRepo{db: db}

			coverInitExpect(mock, tt.state)
			got, err := r.Add(tt.state)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want.id, got[0].ID)
		})
	}
}

func TestDBRepo_Clear(t *testing.T) {
	for _, tt := range getClearTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer db.Close()
			r := DBRepo{db: db}

			q := regexp.QuoteMeta(`INSERT INTO urls(id, url, uid, deleted) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id`)
			mock.ExpectBegin()
			mock.ExpectPrepare(q)
			for _, v := range tt.state {
				mock.ExpectQuery(q).
					WithArgs(v.ID, v.URL, v.UID, v.Deleted).
					WillReturnRows(mock.NewRows([]string{"id"}).AddRow(v.ID))
			}
			mock.ExpectCommit()
			mock.ExpectExec("UPDATE urls SET deleted = true").
				WillReturnResult(sqlmock.NewResult(1, 1))

			if _, err := r.Add(tt.state); err != nil {
				t.Fatal(err)
			}
			r.Clear()
		})
	}
}

func TestDBRepo_Delete(t *testing.T) {
	for _, tt := range getDeleteTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer db.Close()
			r := DBRepo{db: db}

			ids := make([]string, len(tt.state))
			for i, sURL := range tt.state {
				ids[i] = sURL.ID
			}

			coverInitExpect(mock, tt.state)
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE urls SET deleted = true WHERE uid = $1 AND id = any($2)`)).
				WithArgs(UserID, pq.Array(ids)).
				WillReturnResult(sqlmock.NewResult(1, 1))

			if _, err := r.Add(tt.state); err != nil {
				t.Fatal(err)
			}

			err := r.Delete(tt.state)
			log.Error(err)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDBRepo_Get(t *testing.T) {
	for _, tt := range getGetTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer db.Close()
			r := DBRepo{db: db}

			var res ShortURL
			for _, v := range tt.state {
				if v.ID == tt.id {
					res = v
				}
			}

			coverInitExpect(mock, tt.state)
			eq := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM urls WHERE id = $1`)).WithArgs(tt.id)
			if tt.want != "" {
				eq.WillReturnRows(sqlmock.NewRows([]string{"id", "url", "uid", "deleted"}).AddRow(res.ID, res.URL, res.UID, res.Deleted))
			} else {
				eq.WillReturnError(sql.ErrNoRows)
			}

			if _, err := r.Add(tt.state); err != nil {
				t.Fatal(err)
			}
			got, err := r.Get(tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got.URL)
		})
	}
}

func TestDBRepo_GetAll(t *testing.T) {
	for _, tt := range getGetAllTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := getMock(t)
			defer db.Close()
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
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM urls WHERE uid = $1`)).WithArgs(UserID).WillReturnRows(rows)

			if _, err := r.Add(tt.state); err != nil {
				t.Fatal(err)
			}
			got, err := r.GetAll(UserID)
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
			defer db.Close()
			r := DBRepo{db: db}
			var exp int64
			if tt.want {
				exp = 1
			}

			coverInitExpect(mock, tt.state)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM urls WHERE id = $1`)).
				WithArgs(tt.id).
				WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(exp))

			if _, err := r.Add(tt.state); err != nil {
				t.Fatal(err)
			}
			got, err := r.Has(tt.id)
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
	q := regexp.QuoteMeta(`INSERT INTO urls(id, url, uid, deleted) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id`)
	mock.ExpectBegin()
	mock.ExpectPrepare(q)
	for _, v := range state {
		mock.ExpectQuery(q).
			WithArgs(v.ID, v.URL, v.UID, v.Deleted).
			WillReturnRows(mock.NewRows([]string{"id"}).AddRow(v.ID))
	}
	mock.ExpectCommit()
}
