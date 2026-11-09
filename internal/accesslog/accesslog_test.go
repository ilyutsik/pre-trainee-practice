package accesslog

import (
	"os"
	"testing"
)

func TestParseLine(t *testing.T) {
	cases := map[string]struct {
		in   string
		want Line
		ok   bool
	}{
		"empty string":                           {in: "  ", want: Line{}, ok: false},
		"comment string":                         {in: " # this comment  ", want: Line{}, ok: false},
		"invalid string (not 5 params)":          {in: " 203.0.113.10 GET /api/books 200\n", want: Line{}, ok: false},
		"invalid string (path not start with /)": {in: " 203.0.113.10 GET api/books 200 142\n ", want: Line{}, ok: false},
		"invalid string (path has //)":           {in: " 203.0.113.10 GET /api//books 200 142\n ", want: Line{}, ok: false},
		"invalid string (status is not digit)":   {in: " 203.0.113.10 GET /api/books 20F 142\n ", want: Line{}, ok: false},
		"invalid string (bytes is not digit)":    {in: " 203.0.113.10 GET /api/books 200 1P2\n ", want: Line{}, ok: false},
		"valid line": {
			in:   "203.0.113.10 GET /api/books 200 142\n",
			want: Line{IP: "203.0.113.10", Method: "GET", Path: "/api/books", Status: 200, Bytes: 142},
			ok:   true},
		"valid line 2": {
			in:   "198.51.100.2 GET /health 200 2\n",
			want: Line{IP: "198.51.100.2", Method: "GET", Path: "/health", Status: 200, Bytes: 2},
			ok:   true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, ok := ParseLine(c.in)
			if ok != c.ok {
				t.Errorf("%v: ok %v want %v", name, ok, c.ok)
			}
			if !(got == c.want) {
				t.Errorf("%v: got %v want %v", name, got, c.want)
			}
		})
	}
}

func TestSummarize(t *testing.T) {
	text, err := os.ReadFile("../../testdata/access.log")
	if err != nil {
		t.Fatal(err.Error())
	}

	cases := map[string]struct {
		text string
		topN int
		want Summary
	}{
		"empty text": {text: "  ", topN: 3, want: Summary{}},
		"access.log file": {text: string(text),
			topN: 3,
			want: Summary{
				Parsed:    5,
				Skipped:   4,
				UniqueIPs: 3,
				ByStatus:  map[int]int{200: 3, 201: 1, 404: 1},
				TopPaths:  []PathCount{{"/api/books", 3}, {"/api/orders", 1}, {"/health", 1}}}},
		"access.log file (topN > len(topPath))": {text: string(text),
			topN: 10,
			want: Summary{
				Parsed:    5,
				Skipped:   4,
				UniqueIPs: 3,
				ByStatus:  map[int]int{200: 3, 201: 1, 404: 1},
				TopPaths:  []PathCount{{"/api/books", 3}, {"/api/orders", 1}, {"/health", 1}}}},
		"topPaths ordering": {
			text: "198.51.100.2 POST /orders 200 88\n" +
				"198.51.100.2 GET /health 200 2\n" +
				"203.0.113.11 GET /books 200 19",
			topN: 10,
			want: Summary{
				Parsed:    3,
				Skipped:   0,
				UniqueIPs: 2,
				ByStatus:  map[int]int{200: 3},
				TopPaths:  []PathCount{{"/books", 1}, {"/health", 1}, {"/orders", 1}}}},
		"topN < 0": {text: string(text),
			topN: -1,
			want: Summary{}},
		"comment line and empty line": {
			text: "# comment\n   \n   \n",
			topN: 10,
			want: Summary{
				Parsed:    0,
				Skipped:   1,
				UniqueIPs: 0,
				ByStatus:  map[int]int{},
				TopPaths:  []PathCount{}}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := Summarize(c.text, c.topN)
			if !got.Equal(c.want) {
				t.Errorf("%v: got %v want %v", name, got, c.want)
			}
		})
	}
}
