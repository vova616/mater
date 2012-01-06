package collision

import (
	"mater/transform"
	"mater/vect"
	"log"
	"math"
	"strings"
)

type UserData interface{}

type BodyType uint8
const (
	BodyType_Static = iota
	BodyType_Dynamic
)

func (bt BodyType) ToString() string {
	switch bt {
		case BodyType_Static:
			return "Static"
		case BodyType_Dynamic:
			return "Dynamic"
	}
	return "unknown"
}

func (bt *BodyType) FromString(bodyType string) {
	switch strings.ToLower(bodyType) {
		case "static":
			*bt = BodyType_Static
		case "dynamic":
			*bt = BodyType_Dynamic
		default:
			log.Printf("Error: Unknown BodyType \"%v\", BodyType not changed", bodyType)
	}
}

//represents a rigid body
type Body struct {
	//position and rotation of the body
	Transform transform.Transform

	Velocity vect.Vect
	AngularVelocity float64

	Force vect.Vect
	Torque float64

	mass, invMass float64
	i, invI float64

	//all the shapes that make up this body
	Shapes []*Shape

	Space *Space

	Enabled bool
	bodyType BodyType

	IgnoreGravity bool

	//user defined data
	UserData UserData
}

func (body *Body) init() {
	body.Shapes = make([]*Shape, 0, 1)
	body.Enabled = true
}

func NewBody(bodyType BodyType) *Body {
	body := new(Body)
	body.init()
	body.SetBodyType(bodyType)

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

	if shape.ShapeClass == nil {
		log.Printf("Error adding shape: shape.ShapeClass == nil")
		return
	}

	shape.Body = body
	shape.Update()
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
	return body.bodyType == BodyType_Static
}

func (body *Body) SetMass(mass float64) {
	if body.IsStatic() {
		log.Printf("Error: can't change mass of a static body")
		return
	}

	if mass == 0 || math.IsInf(mass, 0) {
		log.Printf("Warning: mass = 0 or mass = inf not valid, setting to 1")
		body.mass = 1
		body.invMass = 1
	} else {
		body.mass = mass
		body.invMass = 1.0 / mass
	}	
}

func (body *Body) Mass() float64 {
	return body.mass
}

func (body *Body) SetInertia(i float64) {
	if body.IsStatic() {
		log.Printf("Error: can't change inertia of a static body")
		return
	}

	body.i = i

	if i == 0 {
		body.invI = 0
	} else {
		body.invI = 1.0 / i
	}	
}

func (body *Body) Inertia() float64 {
	return body.i
}

func (body *Body) UpdateShapes() {
	for _, shape := range body.Shapes {
		shape.Update()
	}
}

func (body *Body) BodyType() BodyType {
	return body.bodyType
}

func (body *Body) SetBodyType(bodyType BodyType) {
	if bodyType == BodyType_Static {
		body.bodyType = BodyType_Static

		body.mass = math.Inf(1)
		body.invMass = 0
		body.i = math.Inf(1)
		body.invI = 0
	} else if bodyType == BodyType_Dynamic {
		body.bodyType = BodyType_Dynamic

		if body.mass == 0.0 || math.IsInf(body.mass, 0) {
			log.Printf("Warning: mass = 0 or mass = inf not valid, setting to 1")
			body.mass = 1
			body.invMass = 1
		}
	} else {
		log.Printf("Error: Unknown BodyType")
	}
}
