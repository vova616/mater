package collision

import (
	"github.com/teomat/mater/aabb"
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
	"log"
)

type shapeProxy struct {
	AABB    aabb.AABB
	ProxyId int
	Shape   *Shape
}

// Datatype used for Shape.CollisionCat and Shape.CollidesWith.
// Change if you need more/less different categories.
type CollisionCategory uint16

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

	proxy shapeProxy

	hash hashValue

	// Surface velocity used when solving for friction.
	Surface_v vect.Vect

	// The categories this shape belongs to.
	// If a shape has any bit in CollidesWith in common with another shape's CollisionCat, then they collide.
	// i.e. if ShapeA.CollidesWith & ShapeB.CollisionCat != 0
	CollisionCat CollisionCategory
	// The categories this shape collides with
	CollidesWith CollisionCategory
}

var shapeIdCounter = hashValue(0)

func newShape() *Shape {
	shape := new(Shape)
	shape.init()
	return shape
}

func (shape *Shape) init() {
	if shape.hash != 0 {
		log.Printf("Warning: Shape already initialized")
		return
	}
	shape.hash = shapeIdCounter
	shapeIdCounter++
	//by default: in group 1, collides only with group 1
	shape.CollisionCat = 1
	shape.CollidesWith = 1
}

// Calls ShapeClass.update and sets the new AABB.
func (shape *Shape) Update() {
	if shape.Body == nil {
		log.Printf("Error: uninitialized shape")
		return
	}
	body := shape.Body

	shape.AABB = shape.ShapeClass.update(body.Transform)
	proxy := &shape.proxy
	proxy.AABB = shape.AABB

	if body.Space != nil {
		d := vect.Sub(body.Transform.Position, body.prevTransform.Position)
		body.Space.BroadPhase.moveProxy(proxy.ProxyId, proxy.AABB, d)
	}
}

func (shape *Shape) createProxy(broadPhase *broadPhase, xf transform.Transform) {
	if shape.proxy.Shape != nil {
		log.Printf("Error: Proxies already created!")
	}

	shape.proxy.Shape = shape
	shape.proxy.AABB = shape.ShapeClass.update(xf)

	shape.proxy.ProxyId = broadPhase.addProxy(shape.proxy)
}

func (shape *Shape) destroyProxy(broadPhase *broadPhase) {
	broadPhase.removeProxy(shape.proxy.ProxyId)
	shape.proxy.Shape = nil
}
