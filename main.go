package main

import (
	"flag"
)

func main() {
	configFlag := initFlag()
	if configFlag.serveMode {
		//serve mode
		listen(configFlag)
	} else if configFlag.clientMode {
		//client mode
		dial(configFlag)
	} else {
		flag.Usage()
	}
}

func initFlag() configFlag {
	var (
		serveMode  bool
		clientMode bool
		address    string
		port       int
	)
	flag.BoolVar(&serveMode, "s", false, "The easyftp serve mode")
	flag.BoolVar(&clientMode, "c", false, "The easyftp client mode")
	flag.StringVar(&address, "a", "127.0.0.1", "The easyftp serve address")
	flag.IntVar(&port, "p", 10021, "The easyftp serve port")
	flag.Parse()
	return configFlag{
		serveMode:  serveMode,
		clientMode: clientMode,
		address:    address,
		port:       port,
	}
}
