package main

import (
	"strings"
)

const (
	echo       = "echo"
	file       = "file"
	fileheader = "fileheader"
	cd         = "cd"
	pwd        = "pwd"
	ls         = "ls"
	lcd        = "lcd"
	lpwd       = "lpwd"
	lls        = "lls"
	get        = "get"
	put        = "put"
	help       = "help"
	quit       = "quit"
)

var supportCmds = "|" + strings.Join([]string{cd, pwd, ls, lcd, lpwd, lls, get, put, help, quit}, "|") + "|"
var localCmds = "|" + strings.Join([]string{lcd, lpwd, lls, help, quit}, "|") + "|"

func isLocalCmd(cmd string) bool {
	return strings.Contains(localCmds, "|"+cmd+"|")
}

func splitCmd(cmd string) (string, string) {
	cmd1, filename := cmd, ""
	if strings.Contains(cmd, ":") {
		ss := strings.Split(cmd, ":")
		cmd1 = ss[0]
		filename = ss[1]
	}
	return cmd1, filename
}
