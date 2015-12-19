package qtype

import (
	"math"
	"testing"
	"testing/quick"
)

func TestVec3Components(t *testing.T) {
	for _, test := range []struct {
		V       Vec3
		X, Y, Z Float
	}{
		{
			V: Vec3{0, 0, 0},
			X: Float(0),
			Y: Float(0),
			Z: Float(0),
		},
		{
			V: Vec3{2, 3, 5},
			X: Float(2),
			Y: Float(3),
			Z: Float(5),
		},
	} {
		if got, want := test.V.X(), test.X; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
		if got, want := test.V.Y(), test.Y; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
		if got, want := test.V.Z(), test.Z; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
	}
}

func TestVec3Len(t *testing.T) {
	for _, test := range []struct {
		V Vec3
		L Float
	}{
		{
			V: Vec3{0, 0, 0},
			L: Float(0),
		},
		{
			V: Vec3{1, 0, 0},
			L: Float(1),
		},
		{
			V: Vec3{2, 3, 5},
			L: Float(6.164414),
		},
	} {
		if got, want := Len(&test.V), test.L; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
	}
}

func TestVec3Normalize(t *testing.T) {
	for _, test := range []struct {
		src  Vec3
		dest Vec3
	}{
		{
			src:  Vec3{0, 0, 0},
			dest: Vec3{0, 0, 0},
		},
		{
			src:  Vec3{1, 1, 1},
			dest: Vec3{0.57735026, 0.57735026, 0.57735026},
		},
	} {
		var out Vec3
		Normalize(&out, &test.src)
		if got, want := out, test.dest; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
	}
}

func TestVec3NormalizeFuzz(t *testing.T) {
	check := func(in Vec3) bool {
		var out Vec3
		l := float64(Normalize(&out, &in))
		if math.IsNaN(l) || math.IsInf(l, 1) || l == 0 {
			return true
		}
		return Len(&out) == 1
	}
	if err := quick.Check(check, nil); err != nil {
		t.Error(err)
	}
}

func TestVec3InvertFuzz(t *testing.T) {
	check := func(in Vec3) bool {
		var out Vec3
		Inverse(&out, &in)
		return out[0] == 1/in[0] && out[1] == 1/in[1] && out[2] == 1/in[2]
	}
	if err := quick.Check(check, nil); err != nil {
		t.Error(err)
	}
}
