package main

import (
	"io"
	"io/ioutil"
	"os"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func deleteFile(path string) {
	os.Remove(path)
}

func tmpFolder() string {
	tmp, err := ioutil.TempDir("", "backup")
	if err != nil {
		panic(err)
	}
	return tmp
}

func getWriter(fname string, fn func(w io.Writer) error) error {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, ConfigPerm)
	if err != nil {
		return err
	}
	defer f.Close()
	return fn(f)
}
