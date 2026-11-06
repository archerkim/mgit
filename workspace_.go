package main

import (
	"os"
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

func listFiles(pathname string) []string {
	IGNORE := []string{".", "..", ".git", ".mgit"}
	entries, err := os.ReadDir(pathname)
	check(err)
	var names []string
	for _, entry := range entries {
		if !contains(IGNORE, entry.Name()) {
			names = append(names, entry.Name())
		}
	}
	return names
}

func readFile(pathname string) []byte {
	bytes, err := os.ReadFile(pathname)
	check(err)
	return bytes
}
