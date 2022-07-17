package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var UserID = "7190e4d4-fd9c-4b"

type AddTestCaseWant struct {
	id  string
	url string
}

type AddTestCase struct {
	name    string
	state   []ShortURL
	want    AddTestCaseWant
	wantErr bool
}

type ClearTestCase struct {
	name  string
	state []ShortURL
}

type GetTestCase struct {
	name    string
	state   []ShortURL
	id      string
	want    string
	wantErr bool
}

type GetAllTestCase struct {
	name    string
	state   []ShortURL
	want    map[string]bool
	wantErr bool
}

type HasTestCase struct {
	name    string
	id      string
	state   []ShortURL
	want    bool
	wantErr bool
}

type DeleteTestCase struct {
	name         string
	state        []ShortURL
	wantDelState bool
	wantErr      bool
}

func getAddTestCases() []AddTestCase {
	return []AddTestCase{
		{
			name:  "Correct URLs",
			state: []ShortURL{{ID: "googl", URL: "https://google.com"}},
			want: AddTestCaseWant{
				id:  "googl",
				url: "https://google.com",
			},
		},
	}
}

func getClearTestCases() []ClearTestCase {
	return []ClearTestCase{
		{
			name:  "Correct clean",
			state: []ShortURL{{ID: "googl", URL: "https://google.com"}},
		},
	}
}

func getGetTestCases() []GetTestCase {
	return []GetTestCase{
		{
			name:    "Missing ID",
			state:   []ShortURL{{ID: "googl", URL: "https://google.com"}},
			id:      "foo",
			wantErr: true,
		},
		{
			name:  "Existing ID",
			state: []ShortURL{{ID: "googl", URL: "https://google.com"}},
			id:    "googl",
			want:  "https://google.com",
		},
	}
}

func getGetAllTestCases() []GetAllTestCase {
	return []GetAllTestCase{
		{
			name:  "One ID present",
			state: []ShortURL{{ID: "googl", URL: "https://google.com", UID: UserID}},
			want:  map[string]bool{"googl": true},
		},
		{
			name: "All IDs present",
			state: []ShortURL{
				{ID: "googl", URL: "https://google.com", UID: UserID},
				{ID: "fcbk", URL: "https://facebook.com", UID: UserID},
			},
			want: map[string]bool{
				"googl": true,
				"fcbk":  true,
			},
		},
		{
			name: "One ID present for user",
			state: []ShortURL{
				{ID: "googl", URL: "https://google.com", UID: UserID},
				{ID: "fcbk", URL: "https://facebook.com", UID: "8201f5e5-ge0d-5c"},
			},
			want: map[string]bool{"googl": true},
		},
	}
}

func getHasTestCases() []HasTestCase {
	return []HasTestCase{
		{
			name:  "Missing ID",
			id:    "foo",
			state: []ShortURL{{ID: "googl", URL: "https://google.com", UID: UserID}},
		},
		{
			name:  "Existing ID",
			id:    "googl",
			state: []ShortURL{{ID: "googl", URL: "https://google.com", UID: UserID}},
			want:  true,
		},
	}
}

func getDeleteTestCases() []DeleteTestCase {
	return []DeleteTestCase{
		{
			name: "Single entry",
			state: []ShortURL{
				{ID: "1", URL: "https://test.com", UID: UserID},
			},
			wantDelState: true,
			wantErr:      false,
		},
		{
			name: "Multiple entries",
			state: []ShortURL{
				{ID: "1", URL: "https://test.com", UID: UserID},
				{ID: "2", URL: "https://test.com", UID: UserID},
			},
			wantDelState: true,
			wantErr:      false,
		},
	}
}

func TestRepo_Add(t *testing.T) {
	for _, tt := range getAddTestCases() {
		for name, r := range getTestRepos(t, "testfile_add") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				t.Parallel()
				got, err := r.Add(tt.state)
				assert.Equal(t, tt.want.id, got[0].ID)
				assert.Equal(t, tt.wantErr, err != nil)
				r.Clear()
			})
		}
	}
}

func TestRepo_Clear(t *testing.T) {
	for _, tt := range getClearTestCases() {
		for name, r := range getTestRepos(t, "testfile_clear") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				t.Parallel()
				if _, err := r.Add(tt.state); err != nil {
					t.Fatal(err)
				}

				r.Clear()

				res, err := r.GetAll(UserID)
				if err != nil {
					t.Error(err)
				}
				assert.Zero(t, len(res))
				r.Clear()
			})
		}
	}
}

func TestRepo_Delete(t *testing.T) {
	t.Parallel()
	for _, tt := range getDeleteTestCases() {
		for name, r := range getTestRepos(t, "testfile_delete") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				if _, err := r.Add(tt.state); err != nil {
					t.Fatal(err)
				}

				err := r.Delete(tt.state)
				assert.Equal(t, tt.wantErr, err != nil)

				for _, sURL := range tt.state {
					stored, err := r.Get(sURL.ID)
					if err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, tt.wantDelState, stored.Deleted)
				}
				r.Clear()
			})
		}
	}
}

func TestRepo_Get(t *testing.T) {
	t.Parallel()
	for _, tt := range getGetTestCases() {
		for name, r := range getTestRepos(t, "testfile_get") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				if _, err := r.Add(tt.state); err != nil {
					t.Fatal(err)
				}

				sURL, err := r.Get(tt.id)
				assert.Equal(t, tt.wantErr, err != nil)
				assert.Equal(t, tt.want, sURL.URL)
			})
			r.Clear()
		}
	}
}

func TestRepo_GetAll(t *testing.T) {
	t.Parallel()
	for _, tt := range getGetAllTestCases() {
		for name, r := range getTestRepos(t, "testfile_getall") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
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
			r.Clear()
		}
	}
}

func TestRepo_Has(t *testing.T) {
	t.Parallel()
	for _, tt := range getHasTestCases() {
		for name, r := range getTestRepos(t, "testfile_has") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				if _, err := r.Add(tt.state); err != nil {
					t.Fatal(err)
				}

				has, err := r.Has(tt.id)
				assert.Equal(t, tt.want, has)
				assert.Equal(t, tt.wantErr, err != nil)
			})
			r.Clear()
		}
	}
}

func getTestName(tname string, rname string) string {
	return rname + "_" + tname
}

func getTestRepos(t *testing.T, fname string) map[string]Storager {
	fr, err := NewFileRepo(fname)
	if err != nil {
		t.Fatal(err)
	}

	return map[string]Storager{
		"memo": NewMemoryRepo(),
		"file": fr,
	}
}
