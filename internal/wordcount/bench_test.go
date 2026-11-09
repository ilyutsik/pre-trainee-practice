package wordcount

import (
	"os"
	"testing"
)

var (
	sinkStr    string
	sinkTokens []string
	sinkCount  map[string]int
)

func loadBench(b *testing.B) string {
	file, err := os.ReadFile("../../testdata/benchFile.txt")
	if err != nil {
		b.Fatal(err)
	}
	return string(file)
}

func BenchmarkTokens(b *testing.B) {
	str := loadBench(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkTokens = Tokens(str)
	}
}

func BenchmarkCount(b *testing.B) {
	str := loadBench(b)
	tokens := Tokens(str)
	for b.Loop() {
		sinkCount = Count(tokens)
	}
}

func BenchmarkJoinPlus(b *testing.B) {
	str := loadBench(b)
	tokens := Tokens(str)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkStr = joinPlus(tokens)
	}
}

func BenchmarkJoinBuilder(b *testing.B) {
	str := loadBench(b)
	tokens := Tokens(str)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkStr = joinBuilder(tokens)
	}
}
