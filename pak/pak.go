package pak

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"unsafe"
)

type ErrNotPak string

func (e ErrNotPak) Error() string { return string(e) }

type ErrNoFile string

func (e ErrNoFile) Error() string { return fmt.Sprintf("pak: no such file %v", e) }

type entity struct{ Offset, Length int }

type Pak struct {
	r   io.ReadSeeker
	dir map[string]*entity
}

func (p *Pak) HasFile(n string) bool {
	_, ok := p.dir[n]
	return ok
}

func (p *Pak) GetFile(n string) ([]byte, error) {
	ent, ok := p.dir[n]
	if !ok {
		return nil, ErrNoFile(n)
	}
	if _, err := p.r.Seek(int64(ent.Offset), 0); err != nil {
		return nil, err
	}
	buf := make([]byte, ent.Length)
	_, err := io.ReadFull(p.r, buf)
	return buf, err
}

func Open(r io.ReadSeeker) (*Pak, error) {
	var hdr struct {
		Id              [4]byte
		DirectoryOffset int32
		DirectoryLength int32
	}
	if err := binary.Read(r, binary.LittleEndian, &hdr); err != nil {
		return nil, err
	}
	var magic = []byte("PACK")
	if !bytes.Equal(magic, hdr.Id[:]) {
		return nil, ErrNotPak(fmt.Sprintf("pak: illegal magic %v", hdr.Id))
	}
	if _, err := r.Seek(int64(hdr.DirectoryOffset), 0); err != nil {
		return nil, err
	}
	var dirent struct {
		Name             [56]byte
		Position, Length int32
	}
	numEntities := int(hdr.DirectoryLength / int32(unsafe.Sizeof(dirent)))
	pak := &Pak{r: r, dir: make(map[string]*entity)}
	for i := 0; i < numEntities; i++ {
		if err := binary.Read(r, binary.LittleEndian, &dirent); err != nil {
			return nil, err
		}
		n := strings.TrimRight(string(dirent.Name[:]), "\x00")
		pak.dir[n] = &entity{
			Offset: int(dirent.Position),
			Length: int(dirent.Length),
		}
	}
	return pak, nil
}
