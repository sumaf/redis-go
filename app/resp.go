package main

import (
	"strconv"
)

type RESP struct {
	Type byte
	Data []string
	Raw []byte
	Count int
}


func AppendPrefix(s []byte, p byte, n int64) []byte {
	s = append(s, p)
	s = strconv.AppendInt(s, n, 10)
	return append(s,'\r','\n')
}

func AppendString (s []byte, simple string) []byte {
	s = append(s, '+') 
	s = append(s, simple...) 
	return append(s, '\r', '\n')
}

func AppendBulkString (s []byte, bulk string) []byte {
	s = AppendPrefix(s, '$', int64(len(bulk)))
	s = append(s,'\r','\n')
	s = append(s, bulk...)
	return append(s,'\r','\n')
}
