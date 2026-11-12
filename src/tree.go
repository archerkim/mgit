package main

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func packTree(entries []entry) []byte {
	sorted := slices.Clone(entries)
	slices.SortFunc(sorted, func(a, b entry) int {
		return strings.Compare(a.name, b.name)
	})

	var buf bytes.Buffer
	for _, e := range sorted {
		buf.WriteString(e.mode)
		buf.WriteString(" ")
		buf.WriteString(e.name)
		buf.WriteByte(0)

		oid, err := hex.DecodeString(e.oid)
		check(err)
		buf.Write(oid)
	}
	return buf.Bytes()
}

func parentDirectories(name string) []string {
	parts := strings.Split(name, "/") // [bin lib jit.rb]
	dirs := make([]string, 0, len(parts)-1)
	for i := 1; i < len(parts); i++ {
		dirs = append(dirs, strings.Join(parts[:i], "/"))
	}
	return dirs // [bin bin/lib]
}

func (db *database) buildTree(root_path string) []entry {
	files := listFiles(root_path)
	var entries []entry
	var path string
	for _, file := range files {
		path = filepath.Join(root_path, file)
		info, err := os.Stat(path)
		check(err)
		parents := parentDirectories(file)
		if len(parents) != 0 {
			data := packTree(db.buildTree(filepath.Join(root_path, parents[0])))
			tree := object{
				oid:   "",
				data:  data,
				_type: "tree",
			}
			db.store(&tree)
			entries = append(entries, entry{oid: tree.oid, name: file, mode: returnMode(info)})
		} else {
			data := readFile(path)
			blob := object{
				oid:   "",
				data:  data,
				_type: "blob",
			}
			db.store(&blob)
			entries = append(entries, entry{oid: blob.oid, name: file, mode: returnMode(info)})
		}
	}
	return entries
}
