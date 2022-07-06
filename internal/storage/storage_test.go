package storage

var UserID = "7190e4d4-fd9c-4b"

type AddTestCaseWant struct {
	id  string
	url string
}

type AddTestCase struct {
	name    string
	batch   []ShortURL
	want    AddTestCaseWant
	wantErr bool
}

type ClearTestCase struct {
	name string
}

type GetTestCase struct {
	name    string
	id      string
	want    string
	wantErr bool
}

type HasTestCase struct {
	name    string
	id      string
	want    bool
	wantErr bool
}

func getAddTestCases() []AddTestCase {
	return []AddTestCase{
		{
			name:  "Correct URLs",
			batch: []ShortURL{{ID: "googl", URL: "https://google.com"}},
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
			name: "Correct clean",
		},
	}
}

func getGetTestCases() []GetTestCase {
	return []GetTestCase{
		{
			name:    "Missing ID",
			id:      "foo",
			wantErr: true,
		},
		{
			name: "Existing ID",
			id:   "googl",
			want: "https://google.com",
		},
	}
}

func getHasTestCases() []HasTestCase {
	return []HasTestCase{
		{
			name: "Missing ID",
			id:   "foo",
		},
		{
			name: "Existing ID",
			id:   "googl",
			want: true,
		},
	}
}
