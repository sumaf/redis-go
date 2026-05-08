package main
import (
	//"github.com/sumaf/redis-go/internal"
	"net"
	"strings"
)



func Dispatch(r RESP, conn net.Conn) {
	switch strings.ToLower(r.Data[0]) {
		case "ping":
			conn.Write([]byte("PONG"))
		case "echo":
			conn.Write([]byte(r.Data[1]))
		default:
			conn.Write([]byte("Unknown command: " + string(r.Data[0])))
	}
}
