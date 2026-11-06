package fileproc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Result struct {
	Path     string
	Lines    int
	Bytes    int
	Err      error // non-nil if the file failed
	Panicked bool  // true if a panic was recovered for this file
}

func Process(paths []string) []Result {
	out := make([]Result, len(paths))

	for idx, path := range paths {
		path = strings.TrimSpace(path)
		out[idx] = processPath(path)
	}
	return out
}

func processPath(path string) (res Result) {
	defer func() {
		if r := recover(); r != nil {
			res.Err = fmt.Errorf("%s", r)
			res.Panicked = true
		}
	}()

	res.Path = path
	file, err := os.Open(path)
	if err != nil {
		res.Err = err
		return res
	}
	defer func() { _ = file.Close() }()

	fileInfo, err := file.Stat()
	if err != nil {
		res.Err = err
		return res
	}
	fileSize := int(fileInfo.Size())
	if fileSize == 0 {
		return res
	}

	data, err := io.ReadAll(file)
	if err != nil {
		res.Err = err
		return res
	}
	var firstLine string
	before, _, ok := bytes.Cut(data, []byte{'\n'})
	if ok {
		firstLine = string(before)
	} else {
		firstLine = string(data)
	}
	if strings.TrimSpace(firstLine) == "PANIC" {
		panic("fileproc: boom")
	}

	lineCount := bytes.Count(data, []byte{'\n'})
	if data[len(data)-1] != '\n' {
		lineCount++
	}
	res.Lines = lineCount
	res.Bytes = fileSize

	return res
}
