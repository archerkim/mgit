package main

import (
	"os"
	"path/filepath"
	"strings"
)

func contains(IGNORE []string, name string) bool {
	for _, ignore := range IGNORE {
		if strings.EqualFold(ignore, name) {
			return true
		}
	}
	return false
}

func listFiles(pathname string, dir ...string) []string {
	IGNORE := []string{".", "..", ".git", ".mgit"}
	var d string
	if len(dir) > 0 {
		d = dir[0]
	} else {
		d = ""
	}
	entries, err := os.ReadDir(pathname)
	check(err)
	var names []string
	for _, entry := range entries {
		if !contains(IGNORE, entry.Name()) {
			if entry.IsDir() {
				names = append(names, listFiles(filepath.Join(pathname, entry.Name()), filepath.Join(d, entry.Name()))...)
			} else {
				names = append(names, filepath.Join(d, entry.Name()))
			}
		}
	}
	return names
}

func readFile(pathname string) []byte {
	bytes, err := os.ReadFile(pathname)
	check(err)
	return bytes
}
