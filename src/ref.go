package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func updateHead(oid string, path string) error {
	l := lockfile{
		filePath: head_path(path),
	}
	isLocked, err := l.hold_for_update()
	if isLocked {
		return fmt.Errorf("Could not acquire lock on the file: %v", err)
	}
	if x := l.write(oid); x != nil {
		return x
	}
	if x := l.write("\n"); x != nil {
		return x
	}
	if x := l.commit(); x != nil {
		return x
	}
	return nil
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
