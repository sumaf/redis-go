package main
import (
	//"github.com/sumaf/redis-go/internal"
	"net"
	"strings"
)



func Dispatch(r RESP, conn net.Conn) {

	if len(r.Data) == 0 {
		conn.Write([]byte("-ERR empty command\r\n"))
		return
	}

	switch strings.ToLower(r.Data[0]) {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			conn.Write(AppendBulkString([]byte{},  r.Data[1]))
		default:
			conn.Write([]byte("Unknown command: " + string(r.Data[0])))
	}
}
