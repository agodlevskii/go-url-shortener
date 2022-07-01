package storage

var UserID = "7190e4d4-fd9c-4b"

type AddTestCaseWant struct {
	id  string
	url string
}

type AddTestCase struct {
	name    string
	repo    Storager
	batch   map[string]string
	want    AddTestCaseWant
	wantErr bool
}

type ClearTestCase struct {
	name string
	repo Storager
}

type GetTestCase struct {
	name    string
	repo    Storager
	id      string
	want    string
	wantErr bool
}

type HasTestCase struct {
	name    string
	repo    Storager
	id      string
	want    bool
	wantErr bool
}

func getAddTestCases(repo Storager) []AddTestCase {
	return []AddTestCase{
		{
			name: "Correct URLs",
			repo: repo,
			batch: map[string]string{
				"googl": "https://google.com",
			},
			want: AddTestCaseWant{
				id:  "googl",
				url: "https://google.com",
			},
		},
	}
}

func getClearTestCases(repo Storager) []ClearTestCase {
	return []ClearTestCase{
		{
			name: "Correct clean",
			repo: repo,
		},
	}
}

func getGetTestCases(repo Storager) []GetTestCase {
	return []GetTestCase{
		{
			name:    "Missing ID",
			repo:    repo,
			id:      "foo",
			wantErr: true,
		},
		{
			name: "Existing ID",
			repo: repo,
			id:   "googl",
			want: "https://google.com",
		},
	}
}

func getHasTestCases(repo Storager) []HasTestCase {
	return []HasTestCase{
		{
			name: "Missing ID",
			repo: repo,
			id:   "foo",
		},
		{
			name: "Existing ID",
			repo: repo,
			id:   "googl",
			want: true,
		},
	}
}
