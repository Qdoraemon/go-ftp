package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func getSubFile(dir string) []string {
	fs := make([]string, 0)
	infos, err := ioutil.ReadDir(dir)
	if err == nil {
		for _, info := range infos {
			f := "f"
			if info.IsDir() {
				f = "d"
			}
			fs = append(fs, fmt.Sprintf("[%s]%s", f, info.Name()))
		}
	}
	return fs
}

func getAbsolutePath(parent, dir string) (string, error) {
	if filepath.IsAbs(dir) {
		_, err := ioutil.ReadDir(dir)
		return dir, err
	}
	dir = strings.Join([]string{parent, dir}, string(filepath.Separator))
	_, err := ioutil.ReadDir(dir)
	return dir, err
}

func getExecDir() string {
	exec, _ := os.Executable()
	return filepath.Dir(exec)
}
