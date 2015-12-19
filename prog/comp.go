package prog

import (
	"encoding/binary"
	"fmt"
	"io"

	. "github.com/matttproud/go-quake/qtype"
)

type Op uint16

const (
	DONE Op = iota
	MUL_F
	MUL_V
	MUL_FV
	MUL_VF
	DIV_F
	ADD_F
	ADD_V
	SUB_F
	SUB_V
	EQ_F
	EQ_V
	EQ_S
	EQ_E
	EQ_FNC
	NE_F
	NE_V
	NE_S
	NE_E
	NE_FNC
	LE
	GE
	LT
	GT
	LOAD_F
	LOAD_V
	LOAD_S
	LOAD_ENT
	LOAD_FLD
	LOAD_FNC
	ADDRESS
	STORE_F
	STORE_V
	STORE_S
	STORE_ENT
	STORE_FLD
	STORE_FNC
	STOREP_F
	STOREP_V
	STOREP_S
	STOREP_ENT
	STOREP_FLD
	STOREP_FNC
	RETURN
	NOT_F
	NOT_V
	NOT_S
	NOT_ENT
	NOT_FNC
	IF
	IFNOT
	CALL0
	CALL1
	CALL2
	CALL3
	CALL4
	CALL5
	CALL6
	CALL7
	CALL8
	STATE
	GOTO
	AND
	OR
	BITAND
	BITOR
	lastOp
)

type Stmt struct {
	// pr_comp.h dstatement_t

	Op                   Op
	First, Second, Third int16
}

const MaxParams = 8

type Func struct {
	// pr_comp.h dstatement_t

	FirstStmt  Int
	ParamStart Int
	Locals     Int
	Profile    Int
	SName      Int
	SFile      Int
	NumParams  Int
	ParamSize  [MaxParams]byte
}

type EType uint16

const (
	ETVoid EType = iota
	ETString
	ETFloat
	ETVector
	ETEntity
	ETField
	ETFunction
	ETPointer
	lastEType

	ETSaveGlobal EType = 1 << 15
)

type Def struct {
	// pr_comp.h ddef_t

	Type   EType
	Offset uint16
	SName  Int
}

type String int32

type stmtInvalidError Stmt

func (s stmtInvalidError) Error() string { return fmt.Sprintf("prog: invalid statement %s", Stmt(s)) }

func decodeStmts(r io.Reader, n int) ([]Stmt, error) {
	stmts := make([]Stmt, n)
	for i := 0; i < n; i++ {
		if err := read(r, &stmts[i]); err != nil {
			return nil, err
		}
		if op := stmts[i].Op; op < 0 || op >= lastOp {
			return nil, stmtInvalidError(stmts[i])
		}
	}
	return stmts, nil
}

func decodeFuncs(r io.Reader, n int) ([]Func, error) {
	funcs := make([]Func, n)
	for i := 0; i < n; i++ {
		if err := read(r, &funcs[i]); err != nil {
			return nil, err
		}
	}
	return funcs, nil
}

type defInvalidError Def

func (d defInvalidError) Error() string { return fmt.Sprintf("prog: invalid definition %s", Def(d)) }

func decodeDefs(r io.Reader, n int) ([]Def, error) {
	defs := make([]Def, n)
	for i := 0; i < n; i++ {
		if err := read(r, &defs[i]); err != nil {
			return nil, err
		}
		if typ := defs[i].Type &^ ETSaveGlobal; typ < 0 || typ >= lastEType {
			return nil, defInvalidError(defs[i])
		}
	}
	return defs, nil
}

type Global Float

func decodeGlobals(r io.Reader) (out []Global, err error) {
	for {
		var g Global
		err = read(r, &g)
		switch err {
		case nil:
			out = append(out, g)
		case io.EOF:
			return out, nil
		default:
			return nil, err
		}
	}
	return out, nil
}

func read(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.LittleEndian, data)
}
