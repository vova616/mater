package vect

import (
	"math"
)

//basic 2d vector.
type Vect struct {
	X, Y float64
}

//adds v2 to the given vector.
func (v1 *Vect) Add(v2 Vect) {
	v1.X += v2.X
	v1.Y += v2.Y
}

//subtracts v2 rom the given vector.
func (v1 *Vect) Sub(v2 Vect) {
	v1.X -= v2.X
	v1.Y -= v2.Y
}

//returns the squared length of the vector.
func (v Vect) LengthSqr() float64 {
	//length of a vector: distance to origin
	return DistSqr(v, Vect{})
}

//returns the length of the vector.
func (v Vect) Length() float64 {
	//length of a vector: distance to origin
	return Dist(v, Vect{})
}

//multiplies the vector by the scalar.
func (v *Vect) Mult(s float64) {
	v.X *= s
	v.Y *= s
}

//normalizes the vector to a length of 1.
func (v *Vect) Normalize() {
	f := 1.0 / v.Length()
	v.X *= f
	v.Y *= f
}

//compare two vectors by value.
func Equals(v1, v2 Vect) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

//adds the input vectors and returns the result.
func Add(v1, v2 Vect) Vect {
	return Vect{v1.X + v2.X, v1.Y + v2.Y}
}

//subtracts the input vectors and returns the result.
func Sub(v1, v2 Vect) Vect {
	return Vect{v1.X - v2.X, v1.Y - v2.Y}
}

//multiplies a vector by a scalar and returns the result.
func Mult(v1 Vect, s float64) Vect {
	return Vect{v1.X * s, v1.Y * s}
}

//returns the square distance between two vectors.
func DistSqr(v1, v2 Vect) float64 {
	return (v1.X - v2.X)*(v1.X - v2.X) + (v1.Y - v2.Y)*(v1.Y - v2.Y)
}

//returns the distance between two vectors.
func Dist(v1, v2 Vect) float64 {
	return math.Sqrt(DistSqr(v1, v2))
}

//returns the squared length of the vector.
func LengthSqr(v Vect) float64 {
	//length of a vector: distance to origin
	return DistSqr(v, Vect{})
}

//returns the length of the vector.
func Length(v Vect) float64 {
	//length of a vector: distance to origin
	return Dist(v, Vect{})
}

//returns a new vector with its x/y values set to the smaller one from the two input values.
//e.g. Min({2, 10}, {8, 3}) would return {2, 3}
func Min(v1, v2 Vect) Vect {
	out := Vect{}

	if v1.X < v2.X {
		out.X = v1.X
	} else {
		out.X = v2.X
	}

	if v1.Y < v2.Y {
		out.Y = v1.Y
	} else {
		out.Y = v2.Y
	}

	return out
}

//returns a new vector with its x/y values set to the bigger one from the two input values.
//e.g. Min({2, 10}, {8, 3}) would return {2, 3}
func Max(v1, v2 Vect) Vect {
	out := Vect{}

	if v1.X > v2.X {
		out.X = v1.X
	} else {
		out.X = v2.X
	}

	if v1.Y > v2.Y {
		out.Y = v1.Y
	} else {
		out.Y = v2.Y
	}

	return out
}

//returns the normalized input vector.
func Normalize(v Vect) Vect {
	f := 1.0 / Length(v)
	return Vect{v.X * f, v.Y * f}
}

//dot product between two vectors.
func Dot(v1, v2 Vect) float64 {
	return v1.X * v2.X + v1.Y * v2.Y
}

//same as CrossVV.
func Cross(a, b Vect) float64 {
	return CrossVV(a, b)
}

//cross product of two vectors.
func CrossVV(a, b Vect) float64 {
	return a.X * b.Y - a.Y * b.X
}

//cross product between a vector and a float64.
//result = {s * a.Y, -s * a.X}
func CrossVD(a Vect, s float64) Vect {
	return Vect{s * a.Y, -s * a.X}
}

//cross product between a float64 and a vector.
//Not the same as CrossVD
//result = {-s * a.Y, s * a.X}
func CrossDV(a Vect, s float64) Vect {
	return Vect{-s * a.Y, s * a.X}
}

//linear interpolation between two vectors by the given scalar
func Lerp(v1, v2 Vect, s float64) Vect {
	return Vect{
		v1.X + (v1.X - v2.X) * s,
		v1.Y + (v1.Y - v2.Y) * s,
	}
}
