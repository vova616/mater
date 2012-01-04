package collision

import (
	"mater/vect"
	"mater/aabb"
	"mater/transform"
	"os"
	"log"
)

//common shape data
type Shape struct {
	Body *Body
	Restitution, Friction float64
	AABB aabb.AABB
	//the actual implementation of the shape
	ShapeClass
}

type ShapeType int
const(
	ShapeType_Circle = 0
	ShapeType_Segment = 1
	numShapes = iota
)

type ShapeClass interface{
	ShapeType() ShapeType
	//update the shape with the new transform and compute the AABB
	Update(xf transform.Transform) aabb.AABB
	//return if the given point is located inside the shape
	TestPoint(xf transform.Transform, point vect.Vect) bool

	//
	MarshalShape(shape *Shape) ([]byte, os.Error)
	UnmarshalShape(shape *Shape, data []byte) (os.Error)
}

func (shape *Shape) Update () {
	if shape.Body == nil {
		log.Printf("Error: uninitialized shape")
		return
	}

	shape.AABB = shape.ShapeClass.Update(shape.Body.Transform)
}
