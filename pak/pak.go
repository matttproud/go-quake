package pak

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

type ErrNotPak string

func (e ErrNotPak) Error() string { return string(e) }

type File struct {
	*io.SectionReader
	Name string
	Size int
}

type Pak struct {
	Files []File
}

func Open(r io.ReaderAt) (*Pak, error) {
	var hdr struct {
		Id              [4]byte
		DirectoryOffset int32
		DirectoryLength int32
	}
	hdrReader := io.NewSectionReader(r, 0, int64(unsafe.Sizeof(hdr)))
	if err := binary.Read(hdrReader, binary.LittleEndian, &hdr); err != nil {
		return nil, err
	}
	var magic = []byte("PACK")
	if !bytes.Equal(magic, hdr.Id[:]) {
		return nil, ErrNotPak(fmt.Sprintf("pak: illegal magic %v", hdr.Id))
	}
	var dirent struct {
		Name             [56]byte
		Position, Length int32
	}
	numEntities := int(hdr.DirectoryLength / int32(unsafe.Sizeof(dirent)))
	pak := &Pak{Files: make([]File, numEntities)}
	dir := io.NewSectionReader(r, int64(hdr.DirectoryOffset), int64(hdr.DirectoryLength))
	for i := 0; i < numEntities; i++ {
		if err := binary.Read(dir, binary.LittleEndian, &dirent); err != nil {
			return nil, err
		}
		tr := trimNull(dirent.Name[:])
		n := string(tr)
		pak.Files[i] = File{
			SectionReader: io.NewSectionReader(r, int64(dirent.Position), int64(dirent.Length)),
			Name:          n,
			Size:          int(dirent.Length),
		}
	}
	return pak, nil
}

func trimNull(b []byte) []byte {
	// strings.TrimRight cannot be used due to abnormalities in the canonical
	// pak data.
	for i, v := range b {
		if v == 0 {
			return b[:i]
		}
	}
	return b
}
