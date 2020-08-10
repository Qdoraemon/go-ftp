package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func listen(configFlag configFlag) {
	address := fmt.Sprintf("%s:%d", configFlag.address, configFlag.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("error = %v\n", err)
		return
	}
	fmt.Printf("[%s]Server started\n", address)
	for {
		conn, _ := listener.Accept()

		//channel for write
		wch := make(chan protocol, 1)

		//channel for read
		rch := make(chan protocol, 1)

		//create client
		client := newClient(conn, wch, rch)

		//client handshake
		handshake(client)

		//poll write to client
		client.write()

		//poll read from client
		client.read()

		//poll read handler from client by poll read
		readHandler(client, rch, wch)
	}
}

func handshake(c *client) {
	fmt.Printf("[%v]Remote connected\n", c.remoteAddr)
}

//read from remote client
func readHandler(c *client, rch chan protocol, wch chan protocol) {
	go func() {
		for {
			p := <-rch
			cmd, filename := splitCmd(p.cmd)
			fullFilename := ""
			if len(filename) > 0 {
				fullFilename = fmt.Sprintf("%s%c%s", c.dir, filepath.Separator, filename)
			}
			fmt.Printf("[%v]%v\n", c.remoteAddr, cmd)
			switch cmd {
			case cd:
				cdCmd(c, p, wch)
				break
			case pwd:
				pwdCmd(c, wch)
				break
			case ls:
				lsCmd(c, wch)
				break
			case get:
				getCmd(c, p, wch)
				break
			case fileheader:
				//create file
				fileheaderCmd(fullFilename)
				break
			case file:
				//transfer file
				fileCmd(fullFilename, p.args)
				break
			}
		}
	}()
}

func pwdCmd(c *client, wch chan protocol) {
	wch <- *encodeEcho([]byte(c.dir))
}

func lsCmd(c *client, wch chan protocol) {
	wch <- *encodeEcho([]byte(strings.Join(getSubFile(c.dir), "\n")))
}

func cdCmd(c *client, p protocol, wch chan protocol) {
	dir := p.argsString()
	str := ""
	newDir, err := getAbsolutePath(c.dir, dir)
	if err != nil {
		str = "Dir not exists"
		wch <- *encodeEcho([]byte(str))
		return
	}
	c.dir = newDir
	wch <- *encodeEcho([]byte(c.dir))
}

func getCmd(c *client, p protocol, wch chan protocol) {
	filename := p.argsString()
	fullFilename := fmt.Sprintf("%v%c%v", c.dir, filepath.Separator, filename)
	f, err := os.Open(fullFilename)
	defer f.Close()
	if err != nil {
		wch <- *encodeEcho([]byte("File not exists"))
		return
	}
	_, err = f.Stat()
	if err != nil {
		wch <- *encodeEcho([]byte("File error:" + err.Error()))
		return
	}

	//start file header
	wch <- *encodeFileheader(filename, []byte{})

	buf := make([]byte, 1*1024*1024) //1MB
	for {
		len, err := f.Read(buf)
		if len > 0 {
			wch <- *encodeFile(filename, buf[:len])
		} else if err == io.EOF {
			//sleep 100ms
			time.Sleep(100 * time.Millisecond)
			wch <- *encodeEcho([]byte("transfer done"))
			break
		}
	}
}

func fileheaderCmd(filename string) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func fileCmd(filename string, buf []byte) {
	f, err := os.OpenFile(filename, os.O_APPEND, 0777)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		_, err = f.Write(buf)
		if err != nil {
			fmt.Println(err)
		}
	}
}
