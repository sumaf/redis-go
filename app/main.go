package main

import (
	"fmt"
	"net"
	"os"
	//"bufio"
	//"sync"
	//"github.com/sumaf/redis-go/internal"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	
	// Stop using the bufio.Reader for now and use byte buffers instead
	//reader := bufio.NewReader(conn)
	
	buf := make([]byte,0,4096)

	for {
		tmp := make([]byte, 1024)
		n, err := conn.Read(tmp)
		if err != nil {
			return
		}

		buf = append(buf, tmp[:n]...)


		for {
			resp, consumed, err := Parsing(buf)
			if err != nil {
				if err.Error() == "incomplete" {
					break
				}
				fmt.Println("failed to parse:", err)
				return
			}

			buf = buf[consumed:]
			Dispatch(resp, conn)
		}
	}
}



func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	
	//var wg sync.WaitGroup 

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}


