package main

import (
	"strconv"
)

type RESP struct {
	Type  byte
	Data  []string
	Raw   []byte
	Count int
}

const (
	respstr = '+'
	resperr = '-'
	respint = ':'
	respblk = '$'
	resparr = '*'
	respnil = '_'
)

func AppendPrefix(s []byte, p byte, n int64) []byte {
	s = append(s, p)
	s = strconv.AppendInt(s, n, 10)
	return append(s, '\r', '\n')
}

func GetSimpleString(simple string) []byte {
	s := append([]byte("+"), simple...)
	return append(s, '\r', '\n')
}

func GetBulkString(bulk string) []byte {
	s := AppendPrefix([]byte{}, '$', int64(len(bulk)))
	s = append(s, bulk...)
	return append(s, '\r', '\n')
}

func GetNullBulkString() []byte {
	return AppendPrefix([]byte{}, '$', -1)
}
