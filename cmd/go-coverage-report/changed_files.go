package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
)

func ParseChangedFiles(filename, prefix string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	var files []string
	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	for i, file := range files {
		files[i] = filepath.Join(prefix, file)
	}

	return files, nil
}
