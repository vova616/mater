package collision

import (
	"github.com/teomat/mater/vect"
	"github.com/teomat/mater/aabb"
	"log"
)

// Holds all bodies, the broadphase and the contactmanager.
type Space struct {
	Enabled    bool
	Gravity    vect.Vect
	Bodies     []*Body
	Iterations int
	Callbacks  struct {
		OnCollision   func(arb *Arbiter)
		ShouldCollide func(sA, sB *Shape) bool
	}

	BroadPhase     *broadPhase
	ContactManager *ContactManager
}

func (space *Space) init() {
	space.Bodies = make([]*Body, 0, 16)
	space.Enabled = true

	space.BroadPhase = newBroadPhase()
	space.ContactManager = newContactManager(space)
}

// Creates a new, empty space.
func NewSpace() *Space {
	space := new(Space)
	space.init()

	return space
}

// Adds the given body to the space.
func (space *Space) AddBody(body *Body) {
	if body.Space != nil {
		log.Printf("Error adding body: body.Space != nil")
		return
	}
	body.Space = space
	body.createProxies()
	body.UpdateShapes()
	space.Bodies = append(space.Bodies, body)
}

// Adds the given body from the space.
func (space *Space) RemoveBody(body *Body) {
	bodies := space.Bodies
	for i, b := range bodies {
		if b == body {
			space.Bodies = append(bodies[:i], bodies[i+1:]...)
			body.destroyProxies()
			return
		}
	}
	log.Printf("Warning removing body: body not found!")
}

// Advances the space by the given timestep.
func (space *Space) Step(dt float64) {

	if dt <= 0.0 {
		return
	}

	inv_dt := 1.0 / dt

	cm := space.ContactManager
	//broadphase
	cm.findNewContacts()

	cm.collide()

	//Integrate forces
	for _, body := range space.Bodies {
		body.UpdateShapes()

		if body.IsStatic() {
			continue
		}

		//b.Velocity += dt * (gravity + b.invMass * b.Force)
		newVel := vect.Mult(body.Force, body.invMass)
		if !body.IgnoreGravity {
			newVel.Add(space.Gravity)
		}
		newVel.Mult(dt)
		body.Velocity.Add(newVel)

		body.AngularVelocity += dt * body.invI * body.Torque

		if Settings.AutoClearForces {
			body.Force = vect.Vect{}
			body.Torque = 0.0
		}
	}

	//Perform pre-steps
	for arb := cm.ArbiterList.Arbiter; arb != nil; arb = arb.Next {
		if arb.ShapeA.IsSensor || arb.ShapeB.IsSensor {
			continue
		}
		arb.preStep(inv_dt)
	}

	//Perform Iterations
	for i := 0; i < Settings.Iterations; i++ {
		for arb := cm.ArbiterList.Arbiter; arb != nil; arb = arb.Next {
			if arb.ShapeA.IsSensor || arb.ShapeB.IsSensor {
				continue
			}
			arb.applyImpulse()
		}
	}

	//Integrate velocities
	for _, body := range space.Bodies {
		if body.IsStatic() {
			continue
		}

		body.Transform.Position.Add(vect.Mult(body.Velocity, dt))

		rot := body.Transform.Angle()
		body.Transform.SetAngle(rot + dt*body.AngularVelocity)

		body.UpdateShapes()
	}
}

func (space *Space) GetDynamicTreeNodes() []dynamicTreeNode {
	return space.BroadPhase._tree._nodes
}

// Queryies the dynamic tree and invokes the callback 
// for each shape whose bounding box overlaps with the given aabb.
// If the callback returns false, the query stops searching for new shapes.
func (space *Space) QueryAABB(callback func(*Shape) bool, aabb aabb.AABB) {
	queryFunc := func(proxyId int) bool{
		proxy := space.BroadPhase.getProxy(proxyId)
		return callback(proxy.Shape)
	}
	space.BroadPhase.query(queryFunc, aabb)
}
