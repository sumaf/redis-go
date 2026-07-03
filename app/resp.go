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

func GetSimpleString(simple string) []byte {
	s := append([]byte{respstr}, simple...)
	return append(s, '\r', '\n')
}

func GetBulkString(bulk string) []byte {
	s := strconv.AppendInt([]byte{respblk}, int64(len(bulk)), 10)
	s = append(s, '\r', '\n')
	s = append(s, bulk...)
	return append(s, '\r', '\n')
}

func GetNullBulkString() []byte {
	s := strconv.AppendInt([]byte{respblk}, -1, 10)
	return append(s, '\r', '\n')
}
