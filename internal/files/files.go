package files

import (
	"os"
	"path/filepath"
	"strings"

	"dcastanho.readson/internal/access"
	md "dcastanho.readson/internal/template"
	// "dcastanho.readson/template/expressions"
)

func GetData(expression string) func() ([]byte, md.Getter) {

	path, accessors := access.SplitJSON(expression)

	result := getFiles(path)
	var sub *[][]byte
	keys := access.ConvertKey(accessors)

	i := 0
	j := 0

	return func() ([]byte, md.Getter) {

		if i < len(*result) {
			dat, _ := os.ReadFile((*result)[i])

			isAr, arr := access.IsArray(keys, dat)
			if isAr {
				sub = arr
			}

			if sub != nil {
				dat = (*sub)[j]
				j++
				if j == len(*sub) {
					sub = nil
					j = 0
					i++
				}
				return dat, access.JSONParserGetterWithBase(keys[1:])
			} else {
				i++
				return dat, access.JSONParserGetterWithBase(keys)
			}
		} else {
			return nil, nil
		}
	}
}

func getFiles(expression string) *[]string {

	var result *[]string

	if isDir(expression) {
		result = getDirFiles(expression)
	} else if isPattern(expression) {
		result = getPatternFiles(expression)
	} else { // is normal file
		result = getSingleFile(expression)
	}

	return result
}

// func isFirstArray()

// File accessors

func getDirFiles(dir string) *[]string {
	result := make([]string, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		// TODO handle read dir error
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filename := entry.Name()
			result = append(result, dir+filename)
		}
	}

	return &result
}

func getSingleFile(file string) *[]string {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		// TODO handle file not exists
	}
	if info.IsDir() {
		// TODO handle file is a dir
	}
	return &[]string{file}
}

func getPatternFiles(pattern string) *[]string {
	res, err := filepath.Glob(pattern)
	if err != nil {
		// TODO handle glob error
	}
	return &res
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
