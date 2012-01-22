package collision

import (
	"github.com/teomat/mater/aabb"
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
)

//If Settings.AutoUpdateShapes is not set, call Update on the parent shape for changes to the A, B and Radius to take effect.
type SegmentShape struct {
	//start/end points of the segment.
	A, B vect.Vect
	//radius of the segment.
	Radius float64

	//local normal. Do not touch!
	N vect.Vect
	//transformed normal. Do not touch!
	Tn vect.Vect
	//transformed start/end points. Do not touch!
	Ta, Tb vect.Vect

	//tangents at the start/end when chained with other segments. Do not touch!
	A_tangent, B_tangent vect.Vect
}

// Creates a new SegmentShape with the given points and radius.
func NewSegment(a, b vect.Vect, r float64) *Shape {
	shape := new(Shape)
	shape.ShapeClass = &SegmentShape{
		A:      a,
		B:      b,
		Radius: r,
	}
	return shape
}

// Returns ShapeType_Segment. Needed to implemet the ShapeClass interface.
func (segment *SegmentShape) ShapeType() ShapeType {
	return ShapeType_Segment
}

//Called to update N, Tn, Ta, Tb and the the bounding box.
func (segment *SegmentShape) update(xf transform.Transform) aabb.AABB {
	a := xf.TransformVect(segment.A)
	b := xf.TransformVect(segment.B)
	segment.Ta = a
	segment.Tb = b
	segment.N = vect.Perp(vect.Normalize(vect.Sub(segment.B, segment.A)))
	segment.Tn = xf.RotateVect(segment.N)

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

// Only returns false for now.
func (segment *SegmentShape) TestPoint(point vect.Vect) bool {
	return false
}
