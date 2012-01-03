package collision

import (
	"mater/vect"
	"mater/transform"
	"mater/aabb"
)

type SegmentShape struct {
	A, B vect.Vect
}

func NewSegmentShape(a, b vect.Vect) *Shape {
	shape := new(Shape)
	shape.ShapeClass = &SegmentShape{
		a, b,
	}
	return shape
}

func (segment *SegmentShape) ShapeType() ShapeType {
	return ShapeType_Segment
}

func (segment *SegmentShape) ComputeAABB(xf transform.Transform) aabb.AABB {
	a := xf.TransformVect(segment.A)
	b := xf.TransformVect(segment.B)

	return aabb.AABB{
		vect.Min(a, b),
		vect.Max(a, b),
	}
}

//TODO:
func (segment *SegmentShape) TestPoint(xf transform.Transform, point vect.Vect) bool {
	return false
}
