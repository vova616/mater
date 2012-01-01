package collision

import (
	"mater/vect"
	"mater/aabb"
	"mater/transform"
	"os"
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
	ShapeType_Circle = iota
)

type ShapeClass interface{
	ShapeType() ShapeType
	//compute the AABB
	ComputeAABB(xf transform.Transform) aabb.AABB
	//return if the given point is located inside the shape
	TestPoint(xf transform.Transform, point vect.Vect) bool

	//
	MarshalShape(shape *Shape) ([]byte, os.Error)
	UnmarshalShape(shape *Shape, data []byte) (os.Error)
}
