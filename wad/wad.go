package wad

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ErrNotWad string

func (e ErrNotWad) Error() string { return string(e) }

func Read(r io.Reader) (interface{}, error) {
	var data struct {
		identification  [4]byte
		numLumps        int32
		infoTableOffset int32
	}
	if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
		return nil, err
	}
	fmt.Printf("data: %v %#v\n", data, data)
	return nil, nil
}
