package main 

import (
	"fmt"
	"strconv"
	"errors"
	"bytes"
)

func Parsing(d []byte) (RESP, int, error) {
	
	if len(d) == 0 {
		return RESP{}, 0, errors.New("incomplete")
	}
	control:= d[0]
	switch control {
		case '$': return parseBulk(d);
		case '+': return parseSimpleString(d);
		case '-': //return parseError(d);
		case ':': //return parseLong(d);
		case '*': return parseArray(d);
		case '~': //return parseSet(d);
		case '%': //return parseMap(d);
		case '#': //return parseBool(d);
		case ',': //return parseDouble(d);
		case '_': //return parseNull(d);
		case '(': //return parseBigNumber(d);
		case '=': //return parseVerbatimString(d);
		case '|': //return parseAttributes(d);
		default: return RESP{}, 0, errors.New("Invalid RESP type")
	}
	return RESP{}, 0, errors.New("Invalid RESP type")
}

func parseArray(d []byte) (RESP, int, error) {
	
	// Count, err := getCount(d) retired until i figure out how to read more than one digit
	header := bytes.Index(d, []byte("\r\n"))
	if header == -1 {
		return RESP{}, 0, errors.New("incomplete")
	}

	length, err := strconv.Atoi(string(d[1:header]))
	if err != nil {
		return RESP{}, 0, err
	}
	if length == -1 {
		return RESP{'*', nil, d[:header],-1}, 0, nil
	}

	var Data []string
	total := header + 2
	
	// skip first CRLF
	values := d[header+2:]

	for i := length; i > 0; i-- {
		element, consumed, err := Parsing(values)
		if err != nil {
			fmt.Println("Failed parsing element:", err)
			return RESP{}, 0, err
		}
		Data = append(Data, element.Data...)
		total += consumed
		values = values[consumed:]
	}

	return RESP{'*', Data, d[:total], length}, total, nil
}

func parseSimpleString(d []byte) (RESP, int, error) {
	header := bytes.Index(d,[]byte("\r\n"))
	if header == -1 {
		return RESP{}, 0, errors.New("incomplete")
	}
	return RESP{'+', []string{string(d[1:header])}, d[:header+2], 0}, header+2, nil
}

func parseBulk(d []byte) (RESP, int, error) {
	
	header := bytes.Index(d,[]byte("\r\n"))
	if header == -1 {
		return RESP{}, 0, errors.New("incomplete")
	}

	length, err := strconv.Atoi(string(d[1:header]))
	if err != nil {
		return RESP{}, 0, err
	}
	if length == -1 {
		return RESP{'$', nil, d[:header+2], -1}, header+2, err
	}

	total := header + 2 + length + 2
	if len(d) < total {
		return RESP{}, 0, errors.New("incomplete")
	}

	return RESP{'$', []string{string(d[header+2:header+2+length])}, d[:total], length}, total, err
}

// Returns index of the header end
func getHeader(d []byte) (int, error) {
	header := bytes.Index(d,[]byte("\r\n"))
	if header == -1 {
		return -1, errors.New("incomplete")
	}
	return header, nil
}
