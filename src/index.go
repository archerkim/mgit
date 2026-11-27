package main

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"maps"
	"os"
	"slices"
	"syscall"
)

const MAX_PATH_SIZE = 0xfff
const HEADER_FORMAT = "a4N2"

type index struct {
	pathname string
	entries  map[string]indexEntry
	lockfile lockfile
}

type indexEntry struct {
	ctime, ctime_nsec, mtime, mtime_nsec, size int64
	dev, ino                                   uint64
	mode, uid, gid                             uint32
	oid, path                                  string
	flags                                      uint16
}

func (index *index) add(pathname string, oid string, info os.FileInfo) {
	indexEntry := indexEntry{
		path:       pathname,
		oid:        oid,
		size:       info.Size(),
		mtime:      info.ModTime().Unix(),
		mtime_nsec: int64(info.ModTime().Nanosecond()),
		flags:      uint16(min(MAX_PATH_SIZE, len(pathname))),
	}
	if m := returnMode(info); m == EXECUTABLE_MODE {
		indexEntry.mode = uint32(0100755)
	} else if m == REGULAR_MODE {
		indexEntry.mode = uint32(0100644)
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		indexEntry.ctime = int64(stat.Ctimespec.Sec)
		indexEntry.ctime_nsec = int64(stat.Ctimespec.Nsec)
		indexEntry.dev = uint64(stat.Dev)
		indexEntry.ino = uint64(stat.Ino)
		indexEntry.uid = stat.Uid
		indexEntry.gid = stat.Gid
	}
	index.entries[pathname] = indexEntry
}

func (index *index) writeUpdates() bool {
	b, err := index.lockfile.hold_for_update()
	check(err)
	if !b {
		return false
	}

	data := []byte("DIRC")
	data = binary.BigEndian.AppendUint32(data, 2)
	data = binary.BigEndian.AppendUint32(data, uint32(len(index.entries)))

	for _, path := range slices.Sorted(maps.Keys(index.entries)) {
		data = append(data, index.entries[path].toBytes()...)
	}

	sum := sha1.Sum(data)
	data = append(data, sum[:]...)

	check(index.lockfile.write(string(data)))
	check(index.lockfile.commit())

	return true
}

func (e indexEntry) toBytes() []byte {
	var b []byte
	for _, v := range []uint32{
		uint32(e.ctime), uint32(e.ctime_nsec),
		uint32(e.mtime), uint32(e.mtime_nsec),
		uint32(e.dev), uint32(e.ino),
		e.mode, e.uid, e.gid,
		uint32(e.size),
	} {
		b = binary.BigEndian.AppendUint32(b, v)
	}

	oid, err := hex.DecodeString(e.oid)
	check(err)
	b = append(b, oid...)

	b = binary.BigEndian.AppendUint16(b, e.flags)

	b = append(b, e.path...)
	b = append(b, 0)

	for len(b)%8 != 0 {
		b = append(b, 0)
	}
	return b
}
