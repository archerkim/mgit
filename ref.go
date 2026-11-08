package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

func updateHead(oid string, path string) {
	file, err := os.OpenFile(head_path(path), os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	defer file.Close()

	_, err = file.WriteString(oid)
	check(err)
}

func readHead(path string) string {

	file, err := os.OpenFile(head_path(path), os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		return ""
	} else {
		check(err)
	}
	defer file.Close()

	oid, err := io.ReadAll(file)
	check(err)

	return string(oid)
}

func head_path(path string) string {
	return filepath.Join(path, "HEAD")
}
