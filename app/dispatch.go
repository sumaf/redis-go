package main

import (
	"net"
	"strings"
	"strconv"
	"errors"
	"time"
)

func Dispatch(r RESP, conn net.Conn, s *Store) {

	if len(r.Data) == 0 {
		conn.Write(GetError("Empty Command"))
		return
	}

	switch strings.ToUpper(r.Data[0]) {
	case "PING":
		conn.Write(GetString("PONG"))
	case "ECHO":
		conn.Write(GetBulkString(r.Data[1]))
	case "SET":
		err := cmdSet(r.Data[1:], s)
		if err != nil {
			conn.Write(GetError(err.Error()))
			return
		}
		conn.Write(GetString("OK"))
	case "GET":
		if len(r.Data) != 2 {
			conn.Write(GetError("Wrong Number of Arguments"))
			return
		}
		value, found := s.Get(r.Data[1])
		if !found {
			conn.Write(GetNullBulkString())
			return
		}
		conn.Write(GetBulkString(value))
	default:
		conn.Write(GetError("Unknown Command: " + string(r.Data[0])))
	}
}


func cmdSet(args []string, s *Store) error {
	if len(args) < 2 {
		return errors.New("Wrong Number of Arguments")
	}

	var opts struct {
		key		string
		value	string
		ttl		time.Duration
		ttlSet	bool
	}

	timeUnit := time.Second 
	opts.key, opts.value, args = args[0], args[1], args[2:]
	for len(args) > 0 {
		switch arg := strings.ToUpper(args[0]); arg {
		case "PX", "PXAT":
			timeUnit = time.Millisecond 
			fallthrough
		case "EX", "EXAT":
			if len(args) < 2 {
				return errors.New("Invalid Syntax")
			}

			if opts.ttlSet {
				return errors.New("Invalid Syntax")
			}

			expire, err := strconv.Atoi(args[1])
			if err != nil {
				return errors.New("Invalid timestamp")
			}

			if arg == "EXAT" || arg == "PXAT" {
				opts.ttl = 0
			} else {
				opts.ttl = time.Duration(expire) * timeUnit
			}
			opts.ttlSet = true

			args = args[2:]
			continue
			
			default:
		}
	}
	
	s.Set(opts.key, opts.value, opts.ttl)

	return nil
}
