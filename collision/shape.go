package collision

import (
	"log"
	"github.com/teomat/mater/aabb"
	"github.com/teomat/mater/transform"

	"github.com/teomat/mater/vect"
)

// Base shape data.
// Holds data all shapetypes have in common.
type Shape struct {
	// The body this shape belongs to.
	Body        *Body
	Restitution float64
	Friction    float64
	AABB        aabb.AABB
	// The actual implementation of the shape.
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
	// Update the shape with the new transform and compute the AABB.
	update(xf transform.Transform) aabb.AABB
	// Returns if the given point is located inside the shape.
	TestPoint(point vect.Vect) bool

	marshalShape(shape *Shape) ([]byte, error)
	unmarshalShape(shape *Shape, data []byte) error
}

// Calls ShapeClass.Update and sets the new AABB.
func (shape *Shape) Update() {
	if shape.Body == nil {
		log.Printf("Error: uninitialized shape")
		return
	}

	shape.AABB = shape.ShapeClass.update(shape.Body.Transform)
}
