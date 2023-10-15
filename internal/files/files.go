package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dcastanho.readson/internal/access"
	md "dcastanho.readson/internal/template"
	// "dcastanho.readson/template/expressions"
)

func GetData(expression string) func() *md.ASTContext {

	path, accessors := access.SplitJSON(expression)

	result, err := getFiles(path)

	if err != nil {
		panic(err.Error())
	}

	var sub *[][]byte
	keys := access.ConvertKey(accessors)

	i := 0
	j := 0

	return func() *md.ASTContext {

		if i < len(*result) {

			filename := (*result)[i]

			fmt.Println(filename)
			dat, _ := os.ReadFile(filename)

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
				getter := access.JSONParserGetterWithBase(keys[1:])
				return &md.ASTContext{Data: dat, Getter: getter}
			} else {
				i++
				getter := access.JSONParserGetterWithBase(keys)
				return &md.ASTContext{Data: dat, Getter: getter}
			}
		} else {
			return nil
		}
	}
}

func getFiles(expression string) (*[]string, error) {

	var result *[]string
	var err error

	if isDir(expression) {
		result, err = getDirFiles(expression)
	} else if isPattern(expression) {
		result, err = getPatternFiles(expression)
	} else { // is normal file
		result, err = getSingleFile(expression)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// func isFirstArray()

// File accessors

func getDirFiles(dir string) (*[]string, error) {
	result := make([]string, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			filename := entry.Name()
			result = append(result, dir+filename)
		}
	}

	return &result, nil
}

func getSingleFile(file string) (*[]string, error) {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("File %s does not exist", file))
	}
	if info.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s is a directory does not exist", file))
	}
	return &[]string{file}, nil
}

func getPatternFiles(pattern string) (*[]string, error) {
	res, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return &res, nil
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
