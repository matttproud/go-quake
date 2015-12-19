package qtype

// Origin the the cartesian coordinate for the center.
var Origin = Vec3{0, 0, 0}

// Vec3 is a three-dimensional vector or coordinate value.
type Vec3 [3]Float

// X returns the x-component of the vector.
func (v *Vec3) X() Float { return v[0] }

// Y returns the y-component of the vector.
func (v *Vec3) Y() Float { return v[1] }

// Z returns the z-component of the vector.
func (v *Vec3) Z() Float { return v[2] }

// Len computes the length of the vector.
func Len(v *Vec3) Float {
	var l Float
	l += v[0] * v[0]
	l += v[1] * v[1]
	l += v[2] * v[2]
	return Sqrt(l)
}

// Normalize writes the src vector transformed by its unit length to the
// destination.
func Normalize(dest, src *Vec3) Float {
	l := Len(src)
	if l == 0 {
		return 0
	}
	i := 1 / l
	dest[0] = i * src[0]
	dest[1] = i * src[1]
	dest[2] = i * src[2]
	return l
}

// Inverse writes the src vector transformed by inversion to the destination.
func Inverse(dest, src *Vec3) {
	dest[0] = 1 / src[0]
	dest[1] = 1 / src[1]
	dest[2] = 1 / src[2]
}

// Equal reports whether two vectors are equivalent in value.
func Equal(a, b *Vec3) bool { return *a == *b }

// Scale writes the src vector transformed by the factor to the destination.
func Scale(dest, src *Vec3, factor Float) {
	dest[0] = factor * src[0]
	dest[1] = factor * src[1]
	dest[2] = factor * src[2]
}
