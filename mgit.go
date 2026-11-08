//usr/bin/env go run -C "$(dirname "$0")" . "$@"; exit

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	arg := os.Args[1]
	switch arg {
	case "init":
		fmt.Println("Initializing...")
		var path string
		if len(os.Args) == 2 {
			path = filepath.Dir(".")
		} else if len(os.Args) == 3 {
			path = os.Args[2]
		} else {
			fmt.Println("Invalid arguments")
			os.Exit(1)
		}
		mgit_path := filepath.Join(path, ".mgit")

		for _, dir := range []string{"objects", "refs"} {
			err := os.MkdirAll(filepath.Join(mgit_path, dir), 0755)
			check(err)
		}
		fmt.Println("Initialized empty mgit repository in ", mgit_path)
	case "commit":
		root_path := filepath.Dir(".")
		mgit_path := filepath.Join(root_path, ".mgit")
		db_path := filepath.Join(mgit_path, "objects")
		db := database{pathname: db_path}
		files := listFiles(root_path)
		var entries []entry
		var path string
		for _, file := range files {
			path = filepath.Join(root_path, file)
			data := readFile(path)
			blob := object{
				oid:   "",
				data:  data,
				_type: "blob",
			}
			db.store(&blob)
			fileInfo, err := os.Stat(path)
			check(err)
			entries = append(entries, entry{oid: blob.oid, name: file, mode: returnMode(fileInfo)})
		}
		data := packTree(entries)
		tree := object{
			oid:   "",
			data:  data,
			_type: "tree",
		}
		db.store(&tree)

		parent := readHead(mgit_path)
		name := os.Getenv("GIT_AUTHOR_NAME")
		email := os.Getenv("GIT_AUTHOR_EMAIL")

		author := name + " " + email + " " + time.Now().Format(time.RFC3339)
		fmt.Print("Print commit message (press command + D to finish input): ")
		data, err := io.ReadAll(os.Stdin)
		check(err)
		message := string(data)

		firstLine, _, _ := strings.Cut(message, "\n")
		firstLine = strings.TrimRight(firstLine, "\r")

		commit_data := "tree " + tree.oid + "\n" +
			"author " + author + "\n" +
			"commiter " + author + "\n\n" +
			message + "\n"

		if parent == "" {
			commit_data = "parent " + parent + "\n" + commit_data
		}

		commit := object{
			oid:   "",
			data:  []byte(commit_data),
			_type: "commit",
		}
		db.store(&commit)
		updateHead(commit.oid, mgit_path)

		var is_root string
		if parent == "" {
			is_root = "(root-commit) "
		} else {
			is_root = ""
		}
		fmt.Printf("[%s%s] %s", is_root, commit.oid, firstLine)
	default:
		fmt.Printf("mgit: %s is not an mgit command\n", arg)
	}
}
