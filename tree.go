package main

import (
	"bytes"
	"encoding/hex"
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
