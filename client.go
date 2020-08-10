package main

import (
	"fmt"
	"net"
)

//Define client struct
type client struct {
	dir        string //The client dir
	localAddr  string //The client localAddr
	remoteAddr string //The client remoteAddr
	writer     *proto //The client writer
	reader     *proto //The client reader
}

//New client
func newClient(conn net.Conn, wch chan protocol, rch chan protocol) *client {
	return &client{
		dir:        getExecDir(),
		localAddr:  fmt.Sprintf("%v", conn.LocalAddr()),
		remoteAddr: fmt.Sprintf("%v", conn.RemoteAddr()),
		writer:     newProto(conn, wch),
		reader:     newProto(conn, rch),
	}
}

func (c *client) read() {
	c.reader.read()
}

func (c *client) write() {
	c.writer.write()
}
