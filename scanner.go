package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//Define scanner struct
type scanner struct {
	c      *client       //client
	signal chan bool     //channel for signal
	lch    chan protocol //channel for local write
	wch    chan protocol //channel for write
}

//New scanner
func newScanner(c *client, signal chan bool, lch, wch chan protocol) *scanner {
	return &scanner{
		c:      c,
		signal: signal,
		lch:    lch,
		wch:    wch,
	}
}

//scan cmd
func (s *scanner) scan() {
outer:
	for {
		fmt.Print(">")
		reader := bufio.NewReader(os.Stdin)
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			continue
		}
		p := encode(buf)
		cmd := strings.TrimSpace(p.cmd)
		if len(cmd) <= 0 {
			continue
		}
		if !strings.Contains(supportCmds, fmt.Sprintf("|%s|", cmd)) {
			fmt.Printf("unsupport command : %v\n", cmd)
			continue
		}
		if isLocalCmd(cmd) {
			s.lch <- *p
		} else {
			switch cmd {
			case cd:
				if len(p.argsString()) <= 0 {
					fmt.Println("The dir must be provided")
					continue outer
				}
			case get:
				if len(p.argsString()) <= 0 {
					fmt.Println("The file name must be provided")
					continue outer
				}
			case put:
				putCmd(s.c, s.wch, p)
				continue outer
			}
			s.wch <- *p
		}
		<-s.signal
	}
}

func putCmd(c *client, wch chan protocol, p *protocol) {
	filename := p.argsString()
	if len(filename) <= 0 {
		fmt.Println("The file name must be provided")
		return
	}
	//find local file
	fullFilename := fmt.Sprintf("%s%c%s", c.dir, filepath.Separator, filename)
	f, err := os.Open(fullFilename)
	defer f.Close()
	if err != nil {
		fmt.Println("File not exists")
		return
	}
	_, err = f.Stat()
	if err != nil {
		fmt.Println("File error:" + err.Error())
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
			fmt.Println("transfer done")
			break
		}
	}
}
