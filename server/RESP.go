package server

import (
	"errors"
)

// reads a RESP encoded simple string from data and returns
// the string, the delta, and the error
func readSimpleString(data []byte) (string, int, error) {
	// first character +
	pos := 1
	//loop until we find \r
	for ; data[pos] != '\r'; pos++ {

	}

	return string(data[1:pos]), pos + 2, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)

	}
	return nil, 0, nil
}

func RESPParser(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	value, _, err := DecodeOne(data)
	return value, err
}
