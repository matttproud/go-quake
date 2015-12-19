package main

import (
	"reflect"
	"testing"
)

func TestDecodeCtrl(t *testing.T) {
	for _, test := range []struct {
		data []byte
		ctrl *Ctrl
		err  error
	}{
		{
			data: []byte{},
			ctrl: nil,
			err:  ErrNotCtrl("too short"),
		},
		{
			data: []byte{128, 0, 0, 12, 1, 81, 85, 65, 75, 69, 0, 3},
			ctrl: &Ctrl{Cmd: 1, Data: []byte{81, 85, 65, 75, 69, 0, 3}},
			err:  nil,
		},
		{
			data: []byte{255, 255, 255, 255, 0},
			ctrl: nil,
			err:  ErrNotCtrl("invalid control signature"),
		},
		{
			data: []byte{128, 0, 0, 13, 1, 81, 85, 65, 75, 69, 0, 3},
			ctrl: nil,
			err:  NewErrWrongLen(13, 12),
		},
	} {
		ctrl, err := DecodeCtrl(test.data)
		if got, want := ctrl, test.ctrl; !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want = %v", got, want)
		}
		if got, want := err, test.err; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
	}
}
