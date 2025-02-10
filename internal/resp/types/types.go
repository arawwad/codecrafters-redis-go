package types

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

const (
	StringType    SimpleString = "string"
	NoneType      SimpleString = "none"
	ListType      SimpleString = "list"
	StreamType    SimpleString = "stream"
	NotApplicable SimpleString = ""
)

type RespType interface {
	Marshal() []byte
	Num() (int, bool)
	Type() SimpleString
}

type SimpleString string

func (s SimpleString) Marshal() []byte {
	str := fmt.Sprintf("+%s\r\n", s)
	return []byte(str)
}

func (s SimpleString) Num() (int, bool) {
	if num, err := strconv.Atoi(string(s)); err == nil {
		return num, true
	}
	return 0, false
}

func (SimpleString) Type() SimpleString {
	return StringType
}

type SimpleError string

func (e SimpleError) Marshal() []byte {
	str := fmt.Sprintf("-%s\r\n", e)
	return []byte(str)
}

func (SimpleError) Num() (int, bool) {
	return 0, false
}

func (SimpleError) Type() SimpleString {
	return NotApplicable
}

type Integer int

func (i Integer) Marshal() []byte {
	str := fmt.Sprintf(":%d\r\n", i)
	return []byte(str)
}

func (i Integer) Num() (int, bool) {
	return int(i), true
}

func (Integer) Type() SimpleString {
	return NotApplicable
}

type Boolean bool

func (i Boolean) Marshal() []byte {
	if i {
		return []byte("#t\r\n")
	}
	return []byte("#f\r\n")
}

func (Boolean) Num() (int, bool) {
	return 0, false
}

func (Boolean) Type() SimpleString {
	return NotApplicable
}

type NullBulkString struct{}

func (NullBulkString) Marshal() []byte {
	return []byte("$-1\r\n")
}

func (NullBulkString) Num() (int, bool) {
	return 0, false
}

func (NullBulkString) Type() SimpleString {
	return NotApplicable
}

type BulkString string

func (s BulkString) Marshal() []byte {
	length := utf8.RuneCountInString(string(s))
	str := fmt.Sprintf("$%d\r\n%s\r\n", length, s)

	return []byte(str)
}

func (s BulkString) Num() (int, bool) {
	if num, err := strconv.Atoi(string(s)); err == nil {
		return num, true
	}
	return 0, false
}

func (BulkString) Type() SimpleString {
	return StringType
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

func (Array) Num() (int, bool) {
	return 0, false
}

func (Array) Type() SimpleString {
	return ListType
}

type Stream string

func (s Stream) Marshal() []byte {
	return []byte{}
}

func (Stream) Num() (int, bool) {
	return 0, false
}

func (Stream) Type() SimpleString {
	return StreamType
}
