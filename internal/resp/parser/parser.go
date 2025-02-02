package parser

import (
	"fmt"
	"strconv"
)

type Parser struct {
	Error        error
	input        []byte
	readPosition int
	position     int
}

func New(input []byte) Parser {
	return Parser{
		input:        input,
		readPosition: 1,
		position:     1,
	}
}

func isDigit(input byte) bool {
	return input >= '0' && input <= '9'
}

func (p *Parser) readLength() (int, bool) {
	for isDigit(p.input[p.readPosition]) {
		p.readPosition = p.readPosition + 1
	}

	length, err := strconv.Atoi(string(p.input[p.position:p.readPosition]))
	if err != nil {
		p.Error = err
		return 0, false
	}

	p.position = p.readPosition
	return length, true
}

func (p *Parser) assertSeparator() bool {
	if p.input[p.readPosition] == '\r' && p.input[p.readPosition+1] == '\n' {
		p.readPosition += 2
		p.position += 2

		return true
	}
	p.Error = fmt.Errorf("Invalid format: Expected \\r\\n at position %d", p.readPosition)
	return false
}

func (p *Parser) parseBulkString() (string, bool) {
	p.readPosition += 1
	p.position += 1

	length, ok := p.readLength()
	if !ok {
		return "", false
	}

	ok = p.assertSeparator()
	if !ok {
		return "", false
	}

	str := string(p.input[p.position : p.position+length])
	p.position += length
	p.readPosition += length

	ok = p.assertSeparator()
	if !ok {
		return "", false
	}
	return str, true
}

func (p *Parser) ParseArray() ([]string, bool) {
	length, ok := p.readLength()
	if !ok {
		return nil, false
	}

	result := make([]string, length, length)

	ok = p.assertSeparator()
	if !ok {
		return nil, false
	}

	for i := 0; i < length; i++ {
		str, ok := p.parseBulkString()
		if !ok {
			return nil, false
		}
		result[i] = str
	}

	return result, true
}
