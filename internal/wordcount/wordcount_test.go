package wordcount

import (
	"maps"
	"os"
	"slices"
	"strings"
	"testing"
)

func TestTokens(t *testing.T) {
	s, err := os.ReadFile("../../testdata/words.txt")
	if err != nil {
		t.Fatal(err.Error())
	}
	cases := map[string]struct {
		in   string
		want []string
	}{
		"words.txt": {string(s), []string{"go", "is", "a", "programming", "language", "go", "is", "expressive",
			"concise", "clean", "and", "efficient", "concurrency", "is", "a", "core", "idea", "in", "go",
			"programming", "programming", "in", "go", "is", "fun", "its", "it", "s"}}}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := Tokens(c.in)
			if !slices.Equal(got, c.want) {
				t.Errorf("%v: got %v want %v", name, got, c.want)
			}
		})
	}
}

func TestCount(t *testing.T) {
	cases := map[string]struct {
		in   []string
		want map[string]int
	}{
		"words.txt Tokens": {
			[]string{"go", "is", "a", "programming", "language", "go", "is", "expressive",
				"concise", "clean", "and", "efficient", "concurrency", "is", "a", "core", "idea", "in", "go",
				"programming", "programming", "in", "go", "is", "fun", "its", "it", "s"},
			map[string]int{"a": 2, "and": 1, "clean": 1, "concise": 1, "concurrency": 1, "core": 1, "efficient": 1, "expressive": 1,
				"fun": 1, "go": 4, "idea": 1, "in": 2, "is": 4, "it": 1, "its": 1, "language": 1, "programming": 3, "s": 1}},

		"nil in slice": {nil, map[string]int{}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := Count(c.in)
			if !maps.Equal(got, c.want) {
				t.Errorf("%v: got %#v want %v", name, got, c.want)
			}
		})
	}
}

func TestTopN(t *testing.T) {
	cases := map[string]struct {
		in   map[string]int
		n    int
		want []Pair
	}{
		"nil in": {nil, 0, []Pair{}},
		"n <= 0": {map[string]int{}, -1, []Pair{}},
		"word.txt freq": {map[string]int{"a": 2, "and": 1, "clean": 1, "concise": 1, "concurrency": 1, "core": 1, "efficient": 1, "expressive": 1,
			"fun": 1, "go": 4, "idea": 1, "in": 2, "is": 4, "it": 1, "its": 1, "language": 1, "programming": 3, "s": 1},
			3,
			[]Pair{{"go", 4}, {"is", 4}, {"programming", 3}}},
		"word.txt freq with n > []Pair)": {map[string]int{"a": 2, "and": 1, "clean": 1, "concise": 1, "concurrency": 1, "core": 1, "efficient": 1, "expressive": 1,
			"fun": 1, "go": 4, "idea": 1, "in": 2, "is": 4, "it": 1, "its": 1, "language": 1, "programming": 3, "s": 1},
			52,
			[]Pair{{"go", 4}, {"is", 4}, {"programming", 3}, {"a", 2}, {"in", 2}, {"and", 1},
				{"clean", 1}, {"concise", 1}, {"concurrency", 1}, {"core", 1}, {"efficient", 1},
				{"expressive", 1}, {"fun", 1}, {"idea", 1}, {"it", 1}, {"its", 1}, {"language", 1},
				{"s", 1}}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := TopN(c.in, c.n)
			if !slices.Equal(got, c.want) {
				t.Errorf("%v: got %v want %v", name, got, c.want)
			}
		})
	}
}

func TestDropStop(t *testing.T) {
	cases := map[string]struct {
		tokens []string
		stop   map[string]struct{}
		want   []string
	}{
		"nil tokens":       {tokens: nil, stop: map[string]struct{}{}, want: []string{}},
		"empty tokens":     {tokens: []string{}, stop: map[string]struct{}{"a": {}, "for": {}}, want: []string{}},
		"empty stop":       {tokens: []string{"hello", "my", "name", "is", "Andrei"}, stop: map[string]struct{}{}, want: []string{"hello", "my", "name", "is", "Andrei"}},
		"simple test case": {tokens: []string{"hello", "my", "name", "is", "Andrei"}, stop: map[string]struct{}{"my": {}, "is": {}}, want: []string{"hello", "name", "Andrei"}},
		"big test case": {tokens: []string{"go", "is", "a", "programming", "language", "go", "is", "expressive",
			"concise", "clean", "and", "efficient", "concurrency", "is", "a", "core", "idea", "in", "go",
			"programming", "programming", "in", "go", "is", "fun", "its", "it", "s"},
			stop: map[string]struct{}{"a": {}, "an": {}, "and": {}, "for": {}, "in": {},
				"is": {}, "of": {}, "on": {}, "or": {}, "the": {}, "to": {}},
			want: []string{"go", "programming", "language", "go", "expressive",
				"concise", "clean", "efficient", "concurrency", "core", "idea", "go",
				"programming", "programming", "go", "fun", "its", "it", "s"}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := DropStop(c.tokens, c.stop)
			if !slices.Equal(got, c.want) {
				t.Errorf("%v: got %v want %v", name, got, c.want)
			}
		})
	}
}

func TestStopWords(t *testing.T) {
	s, err := os.ReadFile("../../testdata/stopList.txt")
	if err != nil {
		t.Fatal(err.Error())
	}

	cases := map[string]struct {
		in   string
		want map[string]struct{}
	}{
		"empty stopList.txt": {in: "", want: map[string]struct{}{}},
		"stopList.txt": {in: string(s),
			want: map[string]struct{}{"a": {}, "an": {}, "and": {}, "for": {}, "in": {},
				"is": {}, "of": {}, "on": {}, "or": {}, "the": {}, "to": {}}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := StopWords(c.in)
			if !maps.Equal(got, c.want) {
				t.Errorf("%v: got %v want %v", name, got, c.want)
			}
		})
	}
}

func TestJoinEqual(t *testing.T) {
	cases := map[string]struct {
		in   []string
		want string
	}{
		"nil":         {in: nil, want: ""},
		"empty":       {in: []string{}, want: ""},
		"easy string": {in: []string{"go", "is", "fun", "", "again"}, want: strings.Join([]string{"go", "is", "fun", "", "again"}, "")},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := joinPlus(c.in); got != c.want {
				t.Errorf("joinPlus=%q; want=%q", got, c.want)
			}
			if got := joinBuilder(c.in); got != c.want {
				t.Errorf("joinBuilder=%q; want=%q", got, c.want)
			}
		})
	}
}
