package collision

import (
	"github.com/teomat/mater/vect"
	"github.com/teomat/mater/aabb"
	"github.com/teomat/mater/transform"
)

type CircleShape struct {
	// Center of the circle. Call Update() if changed.
	Position vect.Vect
	// Radius of the circle. Call Update() if changed.
	Radius float64
	// Global center of the circle. Do not touch!
	Tc vect.Vect
}

func NewCircle(pos vect.Vect, radius float64) *Shape {
	shape := new(Shape)
	shape.ShapeClass = &CircleShape{
		Position: pos, 
		Radius: radius,
	}
	return shape
}

func (circle *CircleShape) ShapeType() ShapeType {
	return ShapeType_Circle
}

// Recalculates the global center of the circle and the the bounding box.
func (circle *CircleShape) Update(xf transform.Transform) aabb.AABB {
	//global center of the circle
	center := xf.TransformVect(circle.Position)
	circle.Tc = center
	rv := vect.Vect{circle.Radius, circle.Radius}

	return aabb.AABB{
		vect.Sub(center, rv),
		vect.Add(center, rv),
	}
}

func (circle *CircleShape) TestPoint(point vect.Vect) bool {
	d := vect.Sub(point, circle.Tc)

	return vect.Dot(d, d) <= circle.Radius * circle.Radius
}
