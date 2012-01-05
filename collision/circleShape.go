package collision

import (
	"mater/vect"
	"mater/aabb"
	"mater/transform"
)

//If Settings.AutoUpdateShapes is not set, call Update on the shape for changes to take the Position and Radius to take effect.
//Don't change TC ever.
type CircleShape struct {
	
	Position vect.Vect
	Radius float64
	//Transform center of the circle, exposed for convenience.
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

func (circle *CircleShape) TestPoint(xf transform.Transform, point vect.Vect) bool {
	center := xf.TransformVect(circle.Position)
	d := vect.Sub(point, center)

	return vect.Dot(d, d) <= circle.Radius * circle.Radius
}
