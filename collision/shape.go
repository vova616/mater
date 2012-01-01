package collision

import (
	"mater/vect"
	"mater/aabb"
)

//common shape data
type Shape struct {
	Body *Body
	Restitution, Friction float64
	AABB aabb.AABB
	//the actual iomplementation of the shape
	ShapeClass
}

type ShapeType int
const(
	ShapeType_Circle = iota
)

type ShapeClass interface{
	ShapeType() ShapeType
	//compute the AABB
	ComputeAABB(pos vect.Vect, rot float64) aabb.AABB
	//return if the given point is located inside the shape
	TestPoint(point vect.Vect) bool
}
