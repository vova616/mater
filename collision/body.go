package collision

import (
	"mater/transform"
	"mater/vect"
	"log"
)

type UserData interface{}

//represents a rigid body
type Body struct {
	//position and rotation of the body
	Transform transform.Transform
	//the global position of the body.
	Position vect.Vect
	//the rotation of the body.
	Rotation float64

	Velocity vect.Vect
	AngularVelocity float64

	Force vect.Vect
	Torque vect.Vect

	Friction float64
	mass, invMass float64
	i, invI float64

	//all the shapes that make up this body
	Shapes []*Shape

	Space *Space

	Enabled bool
	isStatic bool

	//user defined data
	UserData UserData
}

func (body *Body) init () {
	body.mass = 1
	body.invMass = 1
	body.Shapes = make([]*Shape, 0, 1)
	body.Enabled = true
}

func NewBody(static bool) *Body {
	body := new(Body)
	body.init()
	body.isStatic = static

	return body
}

//adds the given shape to the body
func (body *Body) AddShape(shape *Shape) {
	if shape == nil {
		log.Printf("Error adding shape: shape == nil")
		return
	}

	if shape.Body != nil {
		log.Printf("Error adding shape: shape.Body != nil")
		return
	}

	if shape.Body == nil {
		log.Printf("Error adding shape: shape.ShapeClass == nil")
		return
	}

	shape.Body = body
	body.Shapes = append(body.Shapes, shape)
}

//removes the given shape from the body
func (body *Body) RemoveShape(shape *Shape) {
	shapes := body.Shapes
	for i, s := range shapes {
		if s == shape {
			body.Shapes = append(shapes[:i],shapes[i+1:]...)
			return
		}
	}
	log.Printf("Warning removing shape: shape not found!")
}

func (body *Body) IsStatic() bool {
	return body.isStatic
}
