package prog

import (
	"bytes"
	"fmt"
	"io"
	"unsafe"
)

type ErrNotProg string

func (e ErrNotProg) Error() string { return string(e) }

type Prog struct {
	CRC16          uint16
	Funcs          []Func
	Strings        *stringRepo
	GlobalDefs     []Def
	FieldDefs      []Def
	Stmts          []Stmt
	Globals        []Global
	GlobalVars     *GlobalVars
	entityDictSize int
}

func Open(r io.Reader) (*Prog, error) {
	// pr_comp.h dprograms_t
	var hdr struct {
		Version int32
		CRC     int32 // of progdefs interface

		StatementsOffset int32
		StatementsCount  int32

		GlobalDefsOffset int32
		GlobalDefsCount  int32

		FieldDefsOffset int32
		FieldDefsCount  int32

		FunctionsOffset int32
		FunctionsCount  int32

		StringsOffset int32
		StringsCount  int32

		GlobalsOffset int32
		GlobalsCount  int32

		EntityFields int32
	}
	crc := newCRC16()
	r = io.TeeReader(r, crc)
	if err := read(r, &hdr); err != nil {
		return nil, err
	}
	if hdr.Version != 6 {
		return nil, ErrNotProg(fmt.Sprintf("unknown version %v", hdr.Version))
	}
	hdrSz := int(unsafe.Sizeof(hdr))
	is := []int{
		int(hdr.StatementsOffset),
		int(hdr.GlobalDefsOffset),
		int(hdr.FieldDefsOffset),
		int(hdr.FunctionsOffset),
		int(hdr.StringsOffset),
		int(hdr.GlobalsOffset)}
	ivs := newInterval(is...)
	var defs bytes.Buffer
	if _, err := defs.ReadFrom(r); err != nil {
		return nil, err
	}
	readerAt := func(i int32) *bytes.Reader {
		start := int(i) - hdrSz
		until := ivs.End(int(i))
		if until != -1 {
			until -= hdrSz
		} else {
			until = defs.Len()
		}
		return bytes.NewReader(defs.Bytes()[start:until])
	}
	bytesAt := func(i int32) []byte {
		start := int(i) - hdrSz
		until := ivs.End(int(i))
		if until != -1 {
			until -= hdrSz
		} else {
			until = defs.Len()
		}
		return defs.Bytes()[start:until]
	}
	funcs, err := decodeFuncs(readerAt(hdr.FunctionsOffset), int(hdr.FunctionsCount))
	if err != nil {
		return nil, err
	}
	strings := newStringRepo(bytesAt(hdr.StringsOffset))
	globalDefs, err := decodeDefs(readerAt(hdr.GlobalDefsOffset), int(hdr.GlobalDefsCount))
	if err != nil {
		return nil, err
	}
	fieldDefs, err := decodeDefs(readerAt(hdr.FieldDefsOffset), int(hdr.FieldDefsCount))
	if err != nil {
		return nil, err
	}
	stmts, err := decodeStmts(readerAt(hdr.StatementsOffset), int(hdr.StatementsCount))
	if err != nil {
		return nil, err
	}
	globalVars, err := decodeGlobalVars(readerAt(hdr.GlobalsOffset))
	if err != nil {
		return nil, err
	}
	globals, err := decodeGlobals(readerAt(hdr.GlobalsOffset))
	if err != nil {
		return nil, err
	}
	prog := &Prog{
		CRC16:          crc.Sum(),
		Funcs:          funcs,
		Strings:        strings,
		GlobalDefs:     globalDefs,
		FieldDefs:      fieldDefs,
		Stmts:          stmts,
		GlobalVars:     globalVars,
		Globals:        globals,
		entityDictSize: int(hdr.EntityFields),
	}
	return prog, nil
}
