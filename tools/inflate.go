package main

import (
	"compress/zlib"
	"io"
	"log"
	"os"
)

func main() {
	r, err := zlib.NewReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}
}
