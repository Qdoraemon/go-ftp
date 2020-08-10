package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const usageString = `cd        into remote dir
pwd       show remote current dir
ls        show files in remote current dir
lcd       into local dir
lpwd      show local current dir
lls       show files in local current dir
get       get files from remote dir
put       put files to remote dir
help      show usage for help
quit      quit client`

var clientDir = ""
var clientSignal = make(chan bool, 1)
var typeCmdSignal = make(chan bool, 1)
var serveReadChannel = make(chan []byte, 1)

func dial(configFlag configFlag) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", configFlag.address, configFlag.port))
	if err != nil {
		fmt.Printf("error = %v\n", err)
		return
	}

	//channel for write
	wch := make(chan protocol, 1)

	//channel for read
	rch := make(chan protocol, 1)

	//channel for local write
	lch := make(chan protocol, 1)

	//channel for signal
	signal := make(chan bool, 1)

	//new client
	client := newClient(conn, wch, rch)

	//connected
	connected(client)

	//poll write to server
	client.write()

	//poll read from server
	client.read()

	//new scanner
	scanner := newScanner(client, signal, lch, wch)

	//poll read local cmd
	readLocalCmd(client, scanner)

	//poll read channel
	readServeHandler(client, signal, rch)

	//poll scanner stdin
	scanner.scan()
}

func readServeHandler(c *client, signal chan bool, rch chan protocol) {
	go func() {
		for {
			p := <-rch
			cmd, filename := splitCmd(p.cmd)
			fullFilename := ""
			if len(filename) > 0 {
				fullFilename = fmt.Sprintf("%s%c%s", c.dir, filepath.Separator, filename)
			}
			switch cmd {
			case echo:
				fmt.Println(p.argsString())
				signal <- true
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

func connected(c *client) {
	fmt.Printf("[%v]Connected success\n", c.remoteAddr)
}

func clientUsage() {
	fmt.Println(usageString)
}

func readLocalCmd(c *client, s *scanner) {
	go func() {
		for {
			p := <-s.lch
			switch p.cmd {
			case quit:
				fmt.Println("Bye")
				os.Exit(0)
				break
			case help:
				clientUsage()
				break
			case lcd:
				lcdCmd(c, p)
				break
			case lpwd:
				fmt.Println(c.dir)
				break
			case lls:
				fmt.Println(strings.Join(getSubFile(c.dir), "\n"))
				break
			}
			s.signal <- true
		}
	}()
}

func lcdCmd(c *client, p protocol) {
	dir := p.argsString()
	if len(dir) <= 0 {
		fmt.Println("The dir must be provided")
		return
	}
	newDir, err := getAbsolutePath(c.dir, dir)
	if err != nil {
		fmt.Println("Dir not exists")
		return
	}
	c.dir = newDir
	fmt.Println(newDir)
}
