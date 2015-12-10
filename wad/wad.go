// Package wad accesses assets stored in WAD2 archives.
package wad

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

type ErrNotWad string

func (e ErrNotWad) Error() string { return string(e) }

type Type int

const (
	TypeNone    Type = 0
	TypeLabel        = 1
	TypeLumpy        = 64
	TypePalette      = 64
	TypeQTex         = 65
	TypeQPic         = 66
	TypeSound        = 67
	TypeMipTex       = 68
)

type Compression int

const (
	CmpNone Compression = 0
	CmpLzss             = 1
)

type File struct {
	*io.SectionReader
	Name        string
	Size        int
	DiskSize    int
	Type        Type
	Compression Compression
}

type Wad struct {
	Files []*File
}

func Open(r io.ReaderAt) (*Wad, error) {
	var hdr struct {
		Id          [4]byte
		LumpCount   int32
		TableOffset int32
	}
	hdrReader := io.NewSectionReader(r, 0, int64(unsafe.Sizeof(hdr)))
	if err := binary.Read(hdrReader, binary.LittleEndian, &hdr); err != nil {
		return nil, err
	}
	var magic = []byte("WAD2")
	if !bytes.Equal(magic, hdr.Id[:]) {
		return nil, ErrNotWad(fmt.Sprintf("wad: illegal magic %v", hdr.Id))
	}
	wad := &Wad{Files: make([]*File, int(hdr.LumpCount))}
	var lump struct {
		Position    int32
		DiskSize    int32
		Size        int32
		Type        byte
		Compression byte
		Pad1, Pad2  byte
		Name        [16]byte
	}
	dir := io.NewSectionReader(r, int64(hdr.TableOffset), int64(unsafe.Sizeof(lump))*int64(hdr.LumpCount))
	for i := 0; i < int(hdr.LumpCount); i++ {
		if err := binary.Read(dir, binary.LittleEndian, &lump); err != nil {
			return nil, err
		}
		tr := cleanName(lump.Name[:])
		n := string(tr)
		wad.Files[i] = &File{
			SectionReader: io.NewSectionReader(r, int64(lump.Position), int64(lump.Size)),
			Name:          n,
			Type:          Type(lump.Type),
			Compression:   Compression(lump.Compression),
		}
	}
	return wad, nil
}

func cleanName(b []byte) []byte {
	var out []byte
	for _, c := range b {
		if c == 0 {
			break
		}
		if c >= 'A' && c <= 'Z' {
			c += ('a' - 'A')
		}
		out = append(out, c)
	}
	return out
}
