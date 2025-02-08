package parser

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
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
		readPosition: 0,
		position:     0,
	}
}

func isDigit(input byte) bool {
	return input >= '0' && input <= '9'
}

func (p *Parser) Parse() (types.RespType, bool) {
	char := p.input[p.position]
	p.position += 1
	p.readPosition += 1
	switch char {
	case '+':
		return p.parseSimpleString()
	case '$':
		return p.parseBulkString()
	case '*':
		return p.parseArray()
	case ':':
		return p.parseInteger()
	case '#':
		return p.parseBoolean()
	default:
		panic(fmt.Sprintf("unrecognized character: %c", char))
	}
}

func (p *Parser) readNumber() (int, bool) {
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

func (p *Parser) parseBoolean() (types.Boolean, bool) {
	val := false
	if p.input[p.readPosition] == 't' {
		val = true
	} else if p.input[p.readPosition] != 'f' {
		panic("Invalid boolean value")
	}
	p.readPosition += 1
	p.position += 1

	ok := p.assertSeparator()
	if !ok {
		return types.Boolean(false), false
	}
	return types.Boolean(val), true
}

func (p *Parser) parseSimpleString() (types.SimpleString, bool) {
	for ; p.input[p.readPosition] == '\r' && p.input[p.readPosition+1] == '\n'; p.readPosition += 1 {
		if p.readPosition > len(p.input) {
			panic("Coudn't find end of simple string")
		}
	}

	str := string(p.input[p.position:p.readPosition])
	p.readPosition += 2
	p.position = p.readPosition

	return types.SimpleString(str), true
}

func (p *Parser) parseInteger() (types.Integer, bool) {
	num, ok := p.readNumber()
	if !ok {
		return 0, false
	}
	return types.Integer(num), true
}

func (p *Parser) parseBulkString() (types.BulkString, bool) {
	length, ok := p.readNumber()
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
	return types.BulkString(str), true
}

func (p *Parser) parseArray() (types.RespType, bool) {
	length, ok := p.readNumber()
	if !ok {
		return nil, false
	}

	result := make(types.Array, length, length)

	ok = p.assertSeparator()
	if !ok {
		return nil, false
	}

	for i := 0; i < length; i++ {
		val, ok := p.Parse()
		if !ok {
			return nil, false
		}
		result[i] = val
	}

	return result, true
}
