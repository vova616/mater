package collision

import (
	"mater/vect"
	"mater/aabb"
	"mater/transform"
)

type CircleShape struct {
	Position vect.Vect
	Radius float64
}

func NewCircle(pos vect.Vect, radius float64) *Shape {
	shape := new(Shape)
	shape.ShapeClass = &CircleShape{
		pos, radius,
	}
	return shape
}

func (circle *CircleShape) ShapeType() ShapeType {
	return ShapeType_Circle
}

func (circle *CircleShape) ComputeAABB(xf transform.Transform) aabb.AABB {
	//global center of the circle
	center := vect.Add(xf.Position, xf.RotateVect(circle.Position))
	rv := vect.Vect{circle.Radius, circle.Radius}

	return aabb.AABB{
		vect.Sub(center, rv),
		vect.Add(center, rv),
	}
}

func (circle *CircleShape) TestPoint(xf transform.Transform, point vect.Vect) bool {
	center := vect.Add(xf.Position, xf.RotateVect(circle.Position))
	d := vect.Sub(point, center)

	return vect.Dot(d, d) <= circle.Radius * circle.Radius
}
