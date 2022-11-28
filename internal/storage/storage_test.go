package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
			state: []ShortURL{{ID: "google", URL: "https://google.com"}},
			want: AddTestCaseWant{
				id:  "google",
				url: "https://google.com",
			},
		},
	}
}

func getClearTestCases() []ClearTestCase {
	return []ClearTestCase{
		{
			name:  "Correct clean",
			state: []ShortURL{{ID: "google", URL: "https://google.com"}},
		},
	}
}

func getGetTestCases() []GetTestCase {
	return []GetTestCase{
		{
			name:    "Missing ID",
			state:   []ShortURL{{ID: "google", URL: "https://google.com"}},
			id:      "foo",
			wantErr: true,
		},
		{
			name:  "Existing ID",
			state: []ShortURL{{ID: "google", URL: "https://google.com"}},
			id:    "google",
			want:  "https://google.com",
		},
	}
}

func getGetAllTestCases() []GetAllTestCase {
	return []GetAllTestCase{
		{
			name:  "One ID present",
			state: []ShortURL{{ID: "google", URL: "https://google.com", UID: UserID}},
			want:  map[string]bool{"google": true},
		},
		{
			name: "All IDs present",
			state: []ShortURL{
				{ID: "google", URL: "https://google.com", UID: UserID},
				{ID: "facebook", URL: "https://facebook.com", UID: UserID},
			},
			want: map[string]bool{
				"google":   true,
				"facebook": true,
			},
		},
		{
			name: "One ID present for user",
			state: []ShortURL{
				{ID: "google", URL: "https://google.com", UID: UserID},
				{ID: "facebook", URL: "https://facebook.com", UID: "8201f5e5-ge0d-5c"},
			},
			want: map[string]bool{"google": true},
		},
	}
}

func getHasTestCases() []HasTestCase {
	return []HasTestCase{
		{
			name:  "Missing ID",
			id:    "foo",
			state: []ShortURL{{ID: "google", URL: "https://google.com", UID: UserID}},
		},
		{
			name:  "Existing ID",
			id:    "google",
			state: []ShortURL{{ID: "google", URL: "https://google.com", UID: UserID}},
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
		for name, r := range getTestRepos(t, "test_file_add") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				t.Parallel()
				got, err := r.Add(context.Background(), tt.state)
				assert.Equal(t, tt.want.id, got[0].ID)
				assert.Equal(t, tt.wantErr, err != nil)
				r.Clear(context.Background())
			})
		}
	}
}

func TestRepo_Clear(t *testing.T) {
	for _, tt := range getClearTestCases() {
		for name, r := range getTestRepos(t, "test_file_clear") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				t.Parallel()
				if _, err := r.Add(context.Background(), tt.state); err != nil {
					t.Fatal(err)
				}

				r.Clear(context.Background())

				res, err := r.GetAll(context.Background(), UserID, false)
				if err != nil {
					t.Error(err)
				}
				assert.Zero(t, len(res))
				r.Clear(context.Background())
			})
		}
	}
}

func TestRepo_Delete(t *testing.T) {
	t.Parallel()
	for _, tt := range getDeleteTestCases() {
		for name, r := range getTestRepos(t, "test_file_delete") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				if _, err := r.Add(context.Background(), tt.state); err != nil {
					t.Fatal(err)
				}

				err := r.Delete(context.Background(), tt.state)
				assert.Equal(t, tt.wantErr, err != nil)

				for _, sURL := range tt.state {
					stored, err := r.Get(context.Background(), sURL.ID)
					if err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, tt.wantDelState, stored.Deleted)
				}
				r.Clear(context.Background())
			})
		}
	}
}

func TestRepo_Get(t *testing.T) {
	t.Parallel()
	for _, tt := range getGetTestCases() {
		for name, r := range getTestRepos(t, "test_file_get") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				if _, err := r.Add(context.Background(), tt.state); err != nil {
					t.Fatal(err)
				}

				sURL, err := r.Get(context.Background(), tt.id)
				assert.Equal(t, tt.wantErr, err != nil)
				assert.Equal(t, tt.want, sURL.URL)
			})
			r.Clear(context.Background())
		}
	}
}

func TestRepo_GetAll(t *testing.T) {
	t.Parallel()
	for _, tt := range getGetAllTestCases() {
		for name, r := range getTestRepos(t, "test_file_get_all") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				if _, err := r.Add(context.Background(), tt.state); err != nil {
					t.Fatal(err)
				}

				got, err := r.GetAll(context.Background(), UserID, false)
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
			r.Clear(context.Background())
		}
	}
}

func TestRepo_Has(t *testing.T) {
	t.Parallel()
	for _, tt := range getHasTestCases() {
		for name, r := range getTestRepos(t, "test_file_has") {
			t.Run(getTestName(tt.name, name), func(t *testing.T) {
				if _, err := r.Add(context.Background(), tt.state); err != nil {
					t.Fatal(err)
				}

				has, err := r.Has(context.Background(), tt.id)
				assert.Equal(t, tt.want, has)
				assert.Equal(t, tt.wantErr, err != nil)
			})
			r.Clear(context.Background())
		}
	}
}

func getTestName(tName string, rName string) string {
	return rName + "_" + tName
}

func getTestRepos(t *testing.T, fName string) map[string]Storager {
	fr, err := NewFileRepo(fName)
	if err != nil {
		t.Fatal(err)
	}

	return map[string]Storager{
		"memo": NewMemoryRepo(),
		"file": fr,
	}
}
