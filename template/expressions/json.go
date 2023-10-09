package expressions

import (
	"strings"

	"github.com/buger/jsonparser"
)

func addNonEmpty(keys []string, sb strings.Builder) []string {
	if sb.Len() > 0 {
		return append(keys, sb.String())
	} else {
		return keys
	}
}

func convertKey(key string) []string {
	const dotByte byte = byte('.')
	const openBracketByte byte = byte('[')
	const closeBracketByte byte = byte(']')

	var keys []string = make([]string, 0)

	sb := strings.Builder{}

	for i := 0; i < len(key); i++ {
		if key[i] == dotByte {
			keys = addNonEmpty(keys, sb)
			sb = strings.Builder{}
		} else if key[i] == openBracketByte {
			keys = addNonEmpty(keys, sb)
			sb = strings.Builder{}
			sb.WriteByte(key[i])
		} else if key[i] == closeBracketByte {
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

func JSONParserGetter(bytes []byte, key string) (string, error) {
	keys := convertKey(key)
	d, _, _, _ := jsonparser.Get(bytes, keys...)
	s, err := jsonparser.ParseString(d)
	return s, err
}
