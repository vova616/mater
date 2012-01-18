package collision

import (
	"github.com/teomat/mater/aabb"
	"github.com/teomat/mater/vect"
	"log"
)

type Space struct {
	Enabled    bool
	Gravity    vect.Vect
	Bodies     []*Body
	Arbiters   []*Arbiter
	Iterations int
	Callbacks  struct {
		OnCollision   func(arb *Arbiter)
		ShouldCollide func(sA, sB *Shape) bool
	}

	BroadPhase *BroadPhase
}

func (space *Space) init() {
	space.Bodies = make([]*Body, 0, 16)
	space.Arbiters = make([]*Arbiter, 0, 32)
	space.Enabled = true

	space.BroadPhase = NewBroadPhase()
}

func NewSpace() *Space {
	space := new(Space)
	space.init()

	return space
}

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

func (space *Space) Step(dt float64) {

	if dt <= 0.0 {
		return
	}

	inv_dt := 1.0 / dt

	//broadphase
	space.Broadphase()

	//Integrate forces
	for _, body := range space.Bodies {
		if Settings.AutoUpdateShapes {
			body.UpdateShapes()
		}
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
	for _, arb := range space.Arbiters {
		if arb.ShapeA.IsSensor || arb.ShapeB.IsSensor {
			continue
		}
		arb.preStep(inv_dt)
	}

	//Perform Iterations
	for i := 0; i < Settings.Iterations; i++ {
		for _, arb := range space.Arbiters {
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

// O(n^2) broad-phase.
// Tries to collide everything with everything else.
func (space *Space) Broadphase() {
	space.Arbiters = make([]*Arbiter, 0, len(space.Arbiters))
	for i := 0; i < len(space.Bodies)-1; i++ {
		bi := space.Bodies[i]

		for j := i + 1; j < len(space.Bodies); j++ {
			bj := space.Bodies[j]

			if bi.IsStatic() && bj.IsStatic() {
				continue
			}

			for _, si := range bi.Shapes {
				for _, sj := range bj.Shapes {

					shouldCollide := space.Callbacks.ShouldCollide
					if shouldCollide != nil && shouldCollide(si, sj) == false {
						continue
					}

					//check aabbs for overlap
					if !aabb.TestOverlap(si.AABB, sj.AABB) {
						continue
					}

					arb := CreateArbiter(si, sj)
					if arb.NumContacts > 0 {
						onCollisionCallback := space.Callbacks.OnCollision
						if onCollisionCallback != nil {
							onCollisionCallback(arb)
						}
						space.Arbiters = append(space.Arbiters, arb)
					}

					/*
						//search if this arbiter already exists
						var oldArb *Arbiter
						index := 0

						for i , arb := range space.Arbiters {
							if arb.Equals(newArb) {
								oldArb = arb
								index = i
								break
							}
						}

						if newArb.NumContacts > 0 {
							//insert or update the arbiter
							if oldArb == nil {
								println(1)
								//insert
								space.Arbiters = append(space.Arbiters, newArb)
							} else {
								println(2)
								//update
								oldArb.Update(newArb.Contacts, newArb.NumContacts)
							}

						} else {
							if oldArb != nil {
								println(3)
								//remove the arbiter
								space.Arbiters = append(space.Arbiters[:index], space.Arbiters[index+1:]...)
							}
							newArb.Delete()
						}
					*/
				}
			}

		}
	}
}

func (space *Space) GetDynamicTreeNodes() []DynamicTreeNode {
	return space.BroadPhase._tree._nodes
}
