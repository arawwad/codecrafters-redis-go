package marshal

import (
	"fmt"
	"unicode/utf8"
)

func MarshalBulkString(input string) []byte {
	length := utf8.RuneCountInString(input)
	str := fmt.Sprintf("$%d\r\n%s\r\n", length, input)

	return []byte(str)
}
