package aabb

import (
	"github.com/teomat/mater/vect"
)

//axis aligned bounding box.
type AABB struct {
	Lower, Upper vect.Vect
}

//returns the center of the aabb
func (aabb *AABB) Center() vect.Vect {
	return vect.Mult(vect.Add(aabb.Lower, aabb.Upper), 0.5)
}

//returns if other is contained inside this aabb.
func (aabb *AABB) Contains(other AABB) bool {
	return aabb.Lower.X <= other.Lower.X &&
		aabb.Upper.X >= other.Upper.X &&
		aabb.Lower.Y <= other.Lower.Y &&
		aabb.Upper.Y >= other.Upper.Y
}

//returns if v is contained inside this aabb.
func (aabb *AABB) ContainsVect(v vect.Vect) bool {
	return aabb.Lower.X <= v.X &&
		aabb.Upper.X >= v.X &&
		aabb.Lower.Y <= v.Y &&
		aabb.Upper.Y >= v.Y
}

//returns an AABB that holds both a and b.
func Combine(a, b AABB) AABB {
	return AABB{
		vect.Min(a.Lower, b.Lower),
		vect.Max(a.Upper, b.Upper),
	}
}

//returns an AABB that holds both a and v.
func Expand(a AABB, v vect.Vect) AABB {
	return AABB{
		vect.Min(a.Lower, v),
		vect.Max(a.Upper, v),
	}
}

//returns the area of the bounding box.
func (aabb *AABB) Area() float64 {
	return (aabb.Upper.X - aabb.Lower.X) * (aabb.Upper.Y - aabb.Lower.Y)
}

func TestOverlap(a, b AABB) bool {
	d1 := vect.Sub(b.Lower, a.Upper)
	d2 := vect.Sub(a.Lower, b.Upper)

	if d1.X > 0.0 || d1.Y > 0.0 {
		return false
	}

	if d2.X > 0.0 || d2.Y > 0.0 {
		return false
	}

	return true
}
