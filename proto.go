package main

import (
	"fmt"
	"net"
	"runtime"
)

//Define proto struct
type proto struct {
	conn net.Conn      //The conn
	cp   chan protocol //channel for protocol
}

func newProto(conn net.Conn, cp chan protocol) *proto {
	return &proto{
		conn: conn,
		cp:   cp,
	}
}

//write to serve
func (pt *proto) write() {
	go func() {
		for {
			p := <-pt.cp
			pt.conn.Write(p.payload())
		}
	}()
}

//read from serve
func (pt *proto) read() {
	go func() {
		buf := make([]byte, 1*1024*1024)
		for {
			len, err := pt.conn.Read(buf)
			if err != nil {
				fmt.Printf("[%v]Remote closed\n", pt.conn.RemoteAddr())
				//exit current goroutine
				runtime.Goexit()
			}
			p := decode(buf[:len])
			pt.cp <- *p
		}
	}()
}
