package main

import (
	"fmt"
	"os"

	"github.com/ilyutsik/pre-trainee-practice/internal/accesslog"
	"github.com/ilyutsik/pre-trainee-practice/internal/pipeline"
	"github.com/ilyutsik/pre-trainee-practice/internal/wordcount"
)

func main() {
	testPipeline()
	fmt.Println("------------------------------------------------------------------")
	testWorldCount()
	fmt.Println("------------------------------------------------------------------")
	testAccessLog()
}

func testAccessLog() {
	text, err := os.ReadFile("testdata/access.log")
	if err != nil {
		fmt.Println(err.Error())
	}
	s := accesslog.Summarize(string(text), 3)
	fmt.Println(s)
}

func testWorldCount() {
	stopWordsString, err := os.ReadFile("testdata/stopList.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	stop := wordcount.StopWords(string(stopWordsString))

	str, err := os.ReadFile("testdata/words.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tokens := wordcount.Tokens(string(str))
	tokensWithoutStopWords := wordcount.DropStop(tokens, stop)
	countTokens := wordcount.Count(tokensWithoutStopWords)
	topNTokens := wordcount.TopN(countTokens, 3)

	for _, v := range topNTokens {
		fmt.Printf("%s %d \n", v.Word, v.Count)
	}
}

func testPipeline() {
	primarySlice := []int{12, 34, 78, 1, 3, 9, 0, -1, -75, 54, 12, 0, -1, 8, 9, -1}
	//var primarySlice []int

	fmt.Printf("primarySlice before %v \n", primarySlice)

	doubleSlice := pipeline.Double(primarySlice)
	filterSlice := pipeline.FilterMin(primarySlice, 12)
	uniqueSlice := pipeline.UniqueStable(primarySlice)

	fmt.Printf("primarySlice after %v \n", primarySlice)
	fmt.Printf("doubleSlice %v \n", doubleSlice)
	fmt.Printf("filterSlice %v \n", filterSlice)
	fmt.Printf("uniqueSlice %v \n", uniqueSlice)
}
