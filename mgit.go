package main

import (
	"fmt"
	"os"
	"path/filepath"
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
		for _, file := range files {
			data := readFile(filepath.Join(root_path, file))
			blob := object{
				oid:   "",
				data:  data,
				_type: "blob",
			}
			db.store(&blob)
		}
	default:
		fmt.Printf("mgit: %s is not an mgit command\n", arg)
	}
}
