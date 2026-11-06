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
			ex, err := os.Executable()
			check(err)
			path = filepath.Dir(ex)
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
	default:
		fmt.Printf("mgit: %s is not an mgit command\n", arg)
	}
}
