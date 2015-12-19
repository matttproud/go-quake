package qtype

import "math"

// Float is a 32-bit floating point value.
type Float float32

// Sqrt computes the square root of the value.
func Sqrt(f Float) Float { return Float(math.Sqrt(float64(f))) }

// Int is a 32-bit signed integer.
type Int int32
