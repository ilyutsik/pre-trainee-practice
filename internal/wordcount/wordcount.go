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
	out := make([]string, 0, len(s))
	tempSlice := make([]rune, 0, 15)

	flush := func() {
		if len(tempSlice) > 0 {
			out = append(out, strings.ToLower(string(tempSlice)))
			tempSlice = tempSlice[:0]
		}
	}

	for _, v := range s {
		if unicode.IsLetter(v) || unicode.IsDigit(v) {
			tempSlice = append(tempSlice, v)
			continue
		}
		flush()
	}
	flush()
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
