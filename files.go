package main

import (
	"os"
	"path/filepath"
	"strings"
)

// TODO consider changing loops to Walks

func GetFiles(expression string) []string {

	var result []string = make([]string, 0)

	if isDir(expression) {
		entries, err := os.ReadDir(expression)
		if err != nil {
			// TODO handle read dir error
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				filename := entry.Name()
				result = append(result, expression+filename)
			}
		}
	} else if isPattern(expression) {
		res, err := filepath.Glob(expression)
		if err != nil {
			// TODO handle glob error
		}
		result = res

	} else { // is normal file
		info, err := os.Stat(expression)
		if os.IsNotExist(err) {
			// TODO handle file not exists
		}
		if info.IsDir() {
			// TODO handle file is a dir
		}
		result = []string{expression}
	}

	return result
}

func isDir(expression string) bool {
	return expression[len(expression)-1] == byte('/')
}

func isPattern(expression string) bool {
	for _, b := range expression {
		if b == '?' || b == '*' || b == '[' || b == ']' {
			return true
		}
	}
	return false
}

func FileName(path string) string {
	dir, name := filepath.Split(path)
	name = strings.Split(name, ".")[0]
	return filepath.Join(dir, name)
}
