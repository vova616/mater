package collision

import (
	"github.com/teomat/mater/aabb"
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
	"log"
)

type ShapeProxy struct {
	AABB aabb.AABB
	ProxyId int
	Shape *Shape
}

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

	// If the shape is a sensor, collisions are reported but not resolved.
	IsSensor bool

	UserData UserData

	proxy ShapeProxy
}

// Calls ShapeClass.Update and sets the new AABB.
func (shape *Shape) Update() {
	if shape.Body == nil {
		log.Printf("Error: uninitialized shape")
		return
	}

	shape.AABB = shape.ShapeClass.update(shape.Body.Transform)
	proxy := &shape.proxy
	proxy.AABB = shape.AABB

	if shape.Body.Space != nil {
		shape.Body.Space.BroadPhase.MoveProxy(proxy.ProxyId, proxy.AABB, vect.Vect{})
	}
}

func (shape *Shape) createProxy(broadPhase *BroadPhase, xf transform.Transform) {
	if shape.proxy.Shape != nil {
		log.Printf("Error: Proxies already created!")
	}

	shape.proxy.Shape = shape
	shape.proxy.AABB = shape.ShapeClass.update(xf)

	shape.proxy.ProxyId = broadPhase.AddProxy(shape.proxy)
}

func (shape *Shape) destroyProxy(broadPhase *BroadPhase) {
	broadPhase.RemoveProxy(shape.proxy.ProxyId)
}
