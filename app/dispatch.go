package main

import (
	//"github.com/sumaf/redis-go/internal"
	"net"
	"strings"
)

func Dispatch(r RESP, conn net.Conn, s *Store) {

	if len(r.Data) == 0 {
		conn.Write([]byte("-ERR empty command\r\n"))
		return
	}

	switch strings.ToLower(r.Data[0]) {
	case "ping":
		conn.Write([]byte("+PONG\r\n"))
	case "echo":
		conn.Write(AppendBulkString([]byte{}, r.Data[1]))
	case "set":
		if len(r.Data) != 3 {
			conn.Write([]byte("-ERR Wrong Number of Arguments\r\n"))
			return
		}
		s.Set(r.Data[1], r.Data[2])
		conn.Write([]byte("+OK\r\n"))
	case "get":
		if len(r.Data) != 2 {
			conn.Write([]byte("-ERR Wrong Number of Arguments\r\n"))
			return
		}
		value, found := s.Get(r.Data[1])
		if !found {
			conn.Write([]byte("$-1\r\n"))
		}
		conn.Write(AppendBulkString([]byte{},value))
	default:
		conn.Write([]byte("-ERR Unknown command: " + string(r.Data[0]) + "\r\n"))
	}
}
