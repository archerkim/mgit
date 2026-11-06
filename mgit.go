package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func contains(IGNORE []string, name string) bool {
	for _, ignore := range IGNORE {
		if strings.EqualFold(ignore, name) {
			return true
		}
	}
	return false
}

func listFiles(pathname string) []string {
	IGNORE := []string{".", "..", ".git"}
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
		//mgit_path := filepath.Join(root_path, ".mgit")
		//db_path := filepath.Join(mgit_path, "objects")

		fmt.Println(listFiles(root_path))
	default:
		fmt.Printf("mgit: %s is not an mgit command\n", arg)
	}
}
