package types

import (
	"bytes"
	"testing"
)

func TestSimpleString(t *testing.T) {
	input := SimpleString("ok")

	if !bytes.Equal(input.Marshal(), []byte("+ok\r\n")) {
		t.Fatalf("unable to marshal simple string")
	}
}

func TestSimpleError(t *testing.T) {
	input := SimpleError("not found")

	if !bytes.Equal(input.Marshal(), []byte("-not found\r\n")) {
		t.Fatalf("unable to marshal simple error")
	}
}

func TestInteger(t *testing.T) {
	input := Integer(123)

	if !bytes.Equal(input.Marshal(), []byte(":123\r\n")) {
		t.Fatalf("unable to marshal integer")
	}
}

func TestBoolean(t *testing.T) {
	input := Boolean(true)

	if !bytes.Equal(input.Marshal(), []byte("#t\r\n")) {
		t.Fatal("unable to marshal boolean")
	}
}

func TestBulkString(t *testing.T) {
	input := BulkString("hello")

	if !bytes.Equal(input.Marshal(), []byte("$5\r\nhello\r\n")) {
		t.Fatal("unable to marshal bulk string")
	}
}

func TestArray(t *testing.T) {
	input := Array{BulkString("hello"), BulkString("world")}

	if !bytes.Equal(input.Marshal(), []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")) {
		t.Fatal("unable to marsha array")
	}
}
