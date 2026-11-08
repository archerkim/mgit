package main

import "os"

const REGULAR_MODE = "0644"
const EXECUTABLE_MODE = "0755"

type entry struct {
	oid  string
	name string
	mode string
}

func returnMode(info os.FileInfo) string {
	mode := info.Mode()
	if mode == 0644 {
		return REGULAR_MODE
	} else if mode == 0755 {
		return EXECUTABLE_MODE
	} else {
		return "undefined"
	}
}
