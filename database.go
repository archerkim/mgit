package main

import (
	"bytes"
	"compress/flate"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
)

var tempChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type database struct {
	pathname string
}

func (database *database) store(object *object) {
	str := object.data
	content := object._type + " " + strconv.Itoa(len(str)) + "\x00" + string(str)

	hashBytes := sha1.Sum([]byte(content))
	object.oid = hex.EncodeToString(hashBytes[:])
	database.writeObject(object.oid, content)
}

func (database *database) writeObject(oid string, content string) {
	path := filepath.Join(database.pathname, oid[:2], oid[2:])
	dirname := filepath.Dir(path)
	temp_path := filepath.Join(dirname, genTempName())

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			check(os.MkdirAll(dirname, 0755))
		} else {
			log.Fatal(err)
		}
	}
	defer file.Close()

	compressed := compressData([]byte(content))
	check(os.WriteFile(temp_path, compressed, 0644))

	check(os.Rename(temp_path, path))
}

func genTempName() string {
	length := 6
	b := make([]byte, length)

	for i := range b {
		b[i] = tempChars[rand.IntN(len(tempChars))]
	}

	return fmt.Sprintf("tmp_obj_%s", string(b))
}

func compressData(content []byte) []byte {
	var buf bytes.Buffer

	writer, err := flate.NewWriter(&buf, flate.BestSpeed)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := writer.Write(content); err != nil {
		log.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
