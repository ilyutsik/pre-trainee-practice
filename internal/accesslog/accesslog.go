package accesslog

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
)

type Line struct {
	IP     string
	Method string
	Path   string
	Status int
	Bytes  int
}

type Summary struct {
	Parsed    int
	Skipped   int
	UniqueIPs int
	ByStatus  map[int]int
	TopPaths  []PathCount
}

type PathCount struct {
	Path  string
	Count int
}

func ParseLine(s string) (Line, bool) {
	s = strings.Trim(s, " ")
	if s == "" {
		return Line{}, false
	}
	if s[0] == '#' {
		return Line{}, false
	}

	params := strings.Split(s, " ")

	if len(params) != 5 {
		return Line{}, false
	}

	ip := params[0]
	method := params[1]
	path := params[2]
	status := params[3]
	bytes := params[4]

	res := Line{}

	res.IP = ip

	res.Method = method

	if !strings.HasPrefix(path, "/") {
		return Line{}, false
	}
	if strings.Contains(path, "//") {
		return Line{}, false
	}
	res.Path = path

	statusInt, err := strconv.Atoi(status)
	if err != nil {
		return Line{}, false
	}
	res.Status = statusInt

	bytesInt, err := strconv.Atoi(bytes)
	if err != nil {
		return Line{}, false
	}
	res.Bytes = bytesInt

	return res, true
}

func Summarize(text string, topN int) Summary {
	text = strings.Trim(text, " ")
	if text == "" || topN <= 0 {
		return Summary{}
	}

	splitText := strings.Split(text, "\n")

	uniqueIPs := make(map[string]struct{})
	summary := Summary{ByStatus: make(map[int]int)}
	topPathMap := make(map[string]int)
	for _, l := range splitText {
		if len(strings.TrimSpace(l)) == 0 {
			continue
		}
		line, ok := ParseLine(l)
		if ok {
			summary.Parsed++
			summary.ByStatus[line.Status]++
			if _, exist := uniqueIPs[line.IP]; !exist {
				summary.UniqueIPs++
				uniqueIPs[line.IP] = struct{}{}
			}
			topPathMap[line.Path]++
			continue
		}
		summary.Skipped++
	}

	topPath := make([]PathCount, 0, len(topPathMap))
	for k, v := range topPathMap {
		topPath = append(topPath, PathCount{Path: k, Count: v})
	}

	slices.SortFunc(topPath, func(a, b PathCount) int {
		if a.Count != b.Count {
			return cmp.Compare(b.Count, a.Count)
		}
		return cmp.Compare(a.Path, b.Path)
	})
	if topN > len(topPath) {
		topN = len(topPath)
	}
	summary.TopPaths = topPath[0:topN]

	return summary
}

func runeSliceToInt(slice []rune) int {
	res := 0
	for _, r := range slice {
		res = res*10 + int(r-'0')
	}
	return res
}
