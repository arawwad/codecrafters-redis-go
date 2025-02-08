package types

import (
	"fmt"
	"unicode/utf8"
)

type RespType interface {
	Marshal() []byte
}

type SimpleString string

func (s SimpleString) Marshal() []byte {
	str := fmt.Sprintf("+%s\r\n", s)
	return []byte(str)
}

type SimpleError string

func (e SimpleError) Marshal() []byte {
	str := fmt.Sprintf("-%s\r\n", e)
	return []byte(str)
}

type Integer int

func (i Integer) Marshal() []byte {
	str := fmt.Sprintf(":%d\r\n", i)
	return []byte(str)
}

type Boolean bool

func (i Boolean) Marshal() []byte {
	if i {
		return []byte("#t\r\n")
	}
	return []byte("#f\r\n")
}

const NullBulkString = "$-1\r\n"

type BulkString string

func (s BulkString) Marshal() []byte {
	length := utf8.RuneCountInString(string(s))
	str := fmt.Sprintf("$%d\r\n%s\r\n", length, s)

	return []byte(str)
}

type Array []RespType

func (a Array) Marshal() []byte {
	length := len(a)
	result := []byte(fmt.Sprintf("*%d\r\n", length))

	for _, v := range a {
		result = append(result, v.Marshal()...)
	}

	return result
}
