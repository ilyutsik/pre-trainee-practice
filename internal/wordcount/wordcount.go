package wordcount

import (
	"cmp"
	"slices"
	"strings"
	"unicode"
)

type Pair struct {
	Word  string
	Count int
}

func Tokens(s string) []string {
	out := make([]string, 0, len(s)/6)
	start := -1
	hasUpper := false

	flush := func(end int) {
		if start < 0 {
			return
		}
		w := s[start:end]
		if hasUpper {
			w = strings.ToLower(w)
		}
		out = append(out, w)
		start = -1
		hasUpper = false
	}

	for i, v := range s {
		if unicode.IsLetter(v) || unicode.IsDigit(v) {
			if start < 0 {
				start = i
			}
			if unicode.IsUpper(v) {
				hasUpper = true
			}
			continue
		}
		flush(i)
	}
	flush(len(s))
	return out
}

func Count(tokens []string) map[string]int {
	if tokens == nil {
		return map[string]int{}
	}
	out := make(map[string]int, len(tokens))
	for _, v := range tokens {
		out[v]++
	}
	return out
}

func TopN(freq map[string]int, n int) []Pair {
	if freq == nil {
		return []Pair{}
	}
	if n <= 0 {
		return []Pair{}
	}
	result := make([]Pair, 0, len(freq))

	for k, v := range freq {
		result = append(result, Pair{k, v})
	}

	slices.SortFunc(result, func(a, b Pair) int {
		if a.Count != b.Count {
			return cmp.Compare(b.Count, a.Count)
		}
		return cmp.Compare(a.Word, b.Word)
	})
	if n > len(result) {
		n = len(result)
	}
	out := make([]Pair, n)
	copy(out, result[0:n])
	return out
}

func DropStop(tokens []string, stop map[string]struct{}) []string {
	if tokens == nil {
		return []string{}
	}
	out := make([]string, 0, len(tokens))
	for _, v := range tokens {
		if _, exist := stop[v]; exist {
			continue
		}
		out = append(out, v)
	}
	return out
}

func StopWords(stopWordsString string) map[string]struct{} {
	word := make([]rune, 0, 15)
	stop := make(map[string]struct{}, len(stopWordsString)/2)

	flush := func() {
		if len(word) != 0 {
			stopWord := strings.ToLower(string(word))
			stop[stopWord] = struct{}{}
			word = word[:0]
		}
	}

	for _, v := range stopWordsString {
		if unicode.IsLetter(v) || unicode.IsDigit(v) {
			word = append(word, v)
			continue
		}
		flush()
	}
	flush()
	return stop
}

func joinPlus(parts []string) string {
	out := ""
	for _, part := range parts {
		out += part
	}
	return out
}

func joinBuilder(parts []string) string {
	var out strings.Builder
	for _, part := range parts {
		out.WriteString(part)
	}
	return out.String()
}
