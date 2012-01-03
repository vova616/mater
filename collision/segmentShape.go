package collision

import (
	"mater/vect"
	"mater/transform"
	"mater/aabb"
)

type SegmentShape struct {
	A, B vect.Vect
	Radius float64

	n, ta, tb, tn vect.Vect
}

func NewSegmentShape(a, b vect.Vect, r float64) *Shape {
	shape := new(Shape)
	shape.ShapeClass = &SegmentShape{
		A: a,
		B: b,
		Radius: r,
	}
	return shape
}

func (segment *SegmentShape) ShapeType() ShapeType {
	return ShapeType_Segment
}

func (segment *SegmentShape) Update(xf transform.Transform) aabb.AABB {
	a := xf.TransformVect(segment.A)
	b := xf.TransformVect(segment.B)
	segment.ta = a
	segment.tb = b
	segment.n = vect.Perp(vect.Normalize(vect.Sub(segment.A, segment.B)))
	segment.tn = xf.RotateVect(segment.n)

	rv := vect.Vect{segment.Radius, segment.Radius}

	min := vect.Min(a, b)
	min.Sub(rv)

	max := vect.Max(a, b)
	max.Add(rv)

	return aabb.AABB{
		min,
		max,
	}
}

//TODO:
func (segment *SegmentShape) TestPoint(xf transform.Transform, point vect.Vect) bool {
	panic("Not yet implemented!")
	return false
}
