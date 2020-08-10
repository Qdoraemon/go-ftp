# easyftp
A easy ftp implementation in Golang

## Usage
```bash
Usage of easyftp:
  -a string
        The easyftp serve address (default "127.0.0.1")
  -c    The easyftp client mode
  -p int
        The easyftp serve port (default 10021)
  -s    The easyftp serve mode
```
## Mode
* listen

```bash
easyftp -s -a ADDRESS -p PORT
```

* dail

```bash
easyftp -c -a ADDRESS -p PORT
```

## Implementation commands

```bash
cd        into remote dir
pwd       show remote current dir
ls        show files in remote current dir
lcd       into local dir
lpwd      show local current dir
lls       show files in local current dir
get       get files from remote dir
put       put files to remote dir
help      show usage for help
quit      quit client
```