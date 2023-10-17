package access

import (
	"strings"

	md "dcastanho.readson/internal/template"

	"github.com/buger/jsonparser"
)

func addNonEmpty(keys []string, sb strings.Builder) []string {
	if sb.Len() > 0 {
		return append(keys, sb.String())
	} else {
		return keys
	}
}

func isArrow(k string, i int) bool {
	if i == len(k)-1 {
		return false
	} else if k[i] == '-' && k[i+1] == '>' {
		return true
	}
	return false
}

func ConvertKey(key string) []string {

	var keys []string = make([]string, 0)

	sb := strings.Builder{}

	for i := 0; i < len(key); i++ {
		if isArrow(key, i) {
			i++
			keys = addNonEmpty(keys, sb)
			sb = strings.Builder{}
		} else if key[i] == '[' {
			keys = addNonEmpty(keys, sb)
			sb = strings.Builder{}
			sb.WriteByte(key[i])
		} else if key[i] == ']' {
			sb.WriteByte(key[i])
			keys = addNonEmpty(keys, sb)
			sb = strings.Builder{}
		} else {
			sb.WriteByte(key[i])
		}
	}

	if sb.Len() > 0 {
		keys = append(keys, sb.String())
	}

	return keys
}

func JSONArrayEach(data []byte, forEach func(curr []byte, t md.ElementType)) error {

	r := func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		tt := fromJSONtoElement(dataType)
		forEach(value, tt)
	}

	_, err := jsonparser.ArrayEach(data, r)
	return err
}

func JSONObjectEach(data []byte, forEach func(prop string, val []byte, dataType md.ElementType)) error {

	r := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		prop := string(key)
		dt := fromJSONtoElement(dataType)
		forEach(prop, value, dt)
		return nil
	}

	err := jsonparser.ObjectEach(data, r)
	return err
}

func fromJSONtoElement(jsonType jsonparser.ValueType) md.ElementType {
	switch jsonType {
	case jsonparser.String:
		return md.String
	case jsonparser.Boolean:
		return md.Boolean
	case jsonparser.Number:
		return md.Number
	case jsonparser.Array:
		return md.Array
	case jsonparser.Object:
		return md.Object
	case jsonparser.Unknown:
	case jsonparser.NotExist:
	case jsonparser.Null:
	default:
		return md.NotExists
	}
	return md.NotExists
}

func JSONParserGetter(bytes []byte, key string) (string, md.ElementType, error) {
	keys := ConvertKey(key)
	d, dtype, _, err := jsonparser.Get(bytes, keys...)

	if err != nil {
		return "", md.NotExists, err
	}

	s, err := jsonparser.ParseString(d)
	return s, fromJSONtoElement(dtype), err
}

func JSONParserGetterWithBase(base []string) func(bytes []byte, key string) (string, md.ElementType, error) {
	return func(bytes []byte, key string) (string, md.ElementType, error) {
		keys := ConvertKey(key)
		keys = append(base, keys...)

		d, dtype, _, err := jsonparser.Get(bytes, keys...)

		if err != nil {
			return "", md.NotExists, err
		}

		s, err := jsonparser.ParseString(d)
		return s, fromJSONtoElement(dtype), err
	}
}

func IsArray(keys []string, bytes []byte) (bool, *[][]byte) {

	dat, tpe, _, _ := jsonparser.Get(bytes, keys...)

	if tpe == jsonparser.Array {
		return true, getArray(dat)
	}
	return false, nil
}

// func ArrayEach(data []byte, cb func(value []byte, dataType jsonparser.ValueType, offset int, err error), keys ...string)

func getArray(data []byte) *[][]byte {
	ar := make([][]byte, 0)

	adder := func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err == nil {
			ar = append(ar, value)
		}
	}

	jsonparser.ArrayEach(data, adder)
	return &ar
}

func SplitJSON(pattern string) (string, string) {
	sb := strings.Builder{}
	found := false
	var actual string
	var accessPattern string

	for i := 0; i < len(pattern); i++ {
		curr := pattern[i]
		if !found && (curr == '[' || isArrow(pattern, i)) {
			actual = sb.String()
			sb = strings.Builder{}
			found = true
		}
		sb.WriteByte(curr)
	}

	if found {
		accessPattern = sb.String()
	} else {
		actual = sb.String()
	}

	return actual, accessPattern
}
