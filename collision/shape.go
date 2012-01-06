package collision

import (
	"log"
	"mater/aabb"
	"mater/transform"

	"mater/vect"
)

//common shape data
type Shape struct {
	Body        *Body
	Restitution float64
	Friction    float64
	AABB        aabb.AABB
	//the actual implementation of the shape
	ShapeClass
}

type ShapeType int

const (
	ShapeType_Circle  = 0
	ShapeType_Segment = 1
	ShapeType_Polygon = 2
	ShapeType_Box     = 3
	numShapes         = iota
)

type ShapeClass interface {
	ShapeType() ShapeType
	//update the shape with the new transform and compute the AABB
	Update(xf transform.Transform) aabb.AABB
	//return if the given point is located inside the shape
	TestPoint(xf transform.Transform, point vect.Vect) bool

	//
	MarshalShape(shape *Shape) ([]byte, error)
	UnmarshalShape(shape *Shape, data []byte) error
}

func (shape *Shape) Update() {
	if shape.Body == nil {
		log.Printf("Error: uninitialized shape")
		return
	}

	shape.AABB = shape.ShapeClass.Update(shape.Body.Transform)
	v := vect.Vect{.1, .1}
	shape.AABB.Lower.Sub(v)
	shape.AABB.Upper.Add(v)
}
