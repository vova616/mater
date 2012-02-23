package collision

import (
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
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

type Body struct {
	// The body's transform.
	// If Settings.AutoUpdateShapes is set to false, you have to call body.UpdateShapes() for the changes to take effect.
	Transform     transform.Transform
	prevTransform transform.Transform

	Velocity        vect.Vect
	AngularVelocity float64

	Force  vect.Vect
	Torque float64

	mass, invMass float64
	i, invI       float64

	// List of shapes that make up this body. Do not touch!
	Shapes []*Shape

	// The space this body is part of. Don not touch!
	Space *Space

	Enabled  bool
	bodyType BodyType

	IgnoreGravity bool

	// Pointer to userdata.
	UserData UserData

	arbiterList *ArbiterEdge

	w_bias float64
	v_bias vect.Vect
}

func (body *Body) init() {
	body.Shapes = make([]*Shape, 0, 1)
	body.Enabled = true
	body.Transform.SetIdentity()
}

// Creates a new body with the given bodytype and mass and inertia set to 1.
func NewBody(bodyType BodyType) *Body {
	body := new(Body)
	body.init()
	body.mass = 1
	body.invMass = 1
	body.i = 1
	body.invI = 1
	body.SetBodyType(bodyType)

	return body
}

// Adds the given shape to the body.
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
	if body.Space != nil {
		shape.createProxy(body.Space.BroadPhase, body.Transform)
	}
	body.Shapes = append(body.Shapes, shape)
}

// Removes the given shape from the body.
func (body *Body) RemoveShape(shape *Shape) {
	if body.Space != nil {
		shape.destroyProxy(body.Space.BroadPhase)
	}

	shapes := body.Shapes
	for i, s := range shapes {
		if s == shape {
			body.Shapes = append(shapes[:i], shapes[i+1:]...)
			return
		}
	}
	log.Printf("Warning removing shape: shape not found!")
}

//called when the body is added to a space to create all shapeProxies
func (body *Body) createProxies() {
	for _, shape := range body.Shapes {
		shape.createProxy(body.Space.BroadPhase, body.Transform)
	}
}

//called when the body is removed from a space to remove all shapeProxies
func (body *Body) destroyProxies() {
	for _, shape := range body.Shapes {
		shape.destroyProxy(body.Space.BroadPhase)
	}
}

// Returns true if body.BodyType is BodyType_Static.
func (body *Body) IsStatic() bool {
	return body.bodyType == BodyType_Static
}

// Set the mass of the body. If the body is Static the mass can not be changed.
// Setting the mass to 0 or Infinity will not have any effect.
func (body *Body) SetMass(mass float64) {
	if body.IsStatic() {
		log.Printf("Error: can't change mass of a static body")
		return
	}

	if mass == 0 || math.IsInf(mass, 0) {
		log.Printf("Warning: mass = 0 or mass = inf not valid, setting to 1")
	} else {
		body.mass = mass
		body.invMass = 1.0 / mass
	}
}

// Returns the body's mass.
func (body *Body) Mass() float64 {
	return body.mass
}

// Sets the inertia of the body. If the body is Static the inertia can not be changed.
// Setting the inertia to 0 or Infinity will both result in the body no longer being able to rotate.
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

// Returns the body's inertia.
func (body *Body) Inertia() float64 {
	return body.i
}

// Updates the body's shapes. Should be called after the body's transform has been changed.
func (body *Body) UpdateShapes() {
	for _, shape := range body.Shapes {
		shape.Update()
	}
	body.prevTransform = body.Transform
}

// Returns the body's type.
func (body *Body) BodyType() BodyType {
	return body.bodyType
}

// Changes the body's type.
// When changed to Static the mass and inertia are set to +Inf.
// When changed to Dynamic and the mass is either 0 or Inf it is reset to 1.
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

func (body *Body) shouldCollide(other *Body) bool {
	if body.bodyType != BodyType_Dynamic &&
		other.bodyType != BodyType_Dynamic {
		return false
	}
	return true
}
