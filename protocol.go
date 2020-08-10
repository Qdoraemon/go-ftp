package main

import (
	"bytes"
	"strings"
)

//Define protocol struct
type protocol struct {
	cmd  string //cmd string
	args []byte //The cmd args
}

//New protocol
func newProtocol(cmd string, args []byte) *protocol {
	if args == nil {
		args = make([]byte, 0)
	}
	return &protocol{
		cmd:  cmd,
		args: args,
	}
}

func encodeEcho(args []byte) *protocol {
	return newProtocol(echo, args)
}

func encodeFileheader(filename string, args []byte) *protocol {
	return newProtocol(fileheader+":"+filename, args)
}

func encodeFile(filename string, args []byte) *protocol {
	return newProtocol(file+":"+filename, args)
}

//Encode protocol from cmdline
func encode(cmdline []byte) *protocol {
	//find space
	cmdline = bytes.TrimSpace(cmdline)
	index := bytes.IndexByte(cmdline, ' ')
	if index == -1 {
		return newProtocol(string(cmdline), nil)
	}
	cmd := string(cmdline[:index])
	return newProtocol(cmd, bytes.TrimSpace(cmdline[index:]))
}

//Decode protocol from buf
func decode(buf []byte) *protocol {
	cmdIndex := bytes.IndexByte(buf, '\n')
	cmd := ""
	args := make([]byte, 0)
	if cmdIndex != -1 {
		cmd = strings.TrimSpace(string(buf[:cmdIndex]))
		args = buf[cmdIndex+1:]
	} else {
		cmd = strings.TrimSpace(string(buf))
	}
	return newProtocol(cmd, args)
}

//The payload to buffer form cmd
func (p *protocol) payload() []byte {
	buf := []byte(p.cmd)
	buf = append(buf, '\n')
	if len(p.args) > 0 {
		buf = append(buf, p.args[0:len(p.args)]...)
	}
	return buf
}

func (p *protocol) argsString() string {
	return strings.TrimSpace(string(p.args))
}
