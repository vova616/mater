package collision

import (
	"mater/vect"
	"mater/aabb"
	"log"
)

type settings struct {
	AccumulateImpulses bool
}

var Settings settings = settings{
	AccumulateImpulses: true,
}

type Space struct {
	Enabled bool
	Gravity vect.Vect
	Bodies []*Body
	Arbiters []*Arbiter
}

func (space *Space) init() {
	space.Bodies = make([]*Body, 0, 16)
	space.Arbiters = make([]*Arbiter, 0, 32)
	space.Enabled = true
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
	space.Bodies = append(space.Bodies, body)
}

func (space *Space) RemoveBody(body *Body) {
	bodies := space.Bodies
	for i, b := range bodies {
		if b == body {
			space.Bodies = append(bodies[:i], bodies[i+1:]...)
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
	_ = inv_dt

	//broadphase
	space.Broadphase()

	//Integrate forces
	for _, body := range space.Bodies {
		if body.IsStatic() {
			continue
		}

		//b.Velocity += dt * (gravity + b.invMass * b.Force)
		newVel := vect.Add(space.Gravity, vect.Mult(body.Force, body.invMass))
		newVel.Mult(dt)
		body.Velocity.Add(newVel)

		body.AngularVelocity += dt * body.invI * body.Torque
	}

	//stuff

	//Integrate velocities
	for _, body := range space.Bodies {
		if body.IsStatic() {
			continue
		}

		body.Transform.Position.Add(vect.Mult(body.Velocity, dt))

		rot := body.Transform.Angle()
		body.Transform.SetAngle(rot + dt * body.AngularVelocity)

		body.UpdateAABBs()
	}
}

func (space *Space) Broadphase() {
	for i := 0; i < len(space.Bodies) - 1; i++ {
		bi := space.Bodies[i]

		for j := i + 1; j < len(space.Bodies); j++ {
			bj := space.Bodies[j]

			if bi.IsStatic() && bj.IsStatic() {
				continue
			}

			for _, si := range bi.Shapes {
				for _, sj := range bj.Shapes {
					//check aabbs for overlap
					if !aabb.TestOverlap(si.AABB, sj.AABB) {
						continue
					}

					newArb := CreateArbiter(si, sj)

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
							//insert
							space.Arbiters = append(space.Arbiters, newArb)
						} else {
							//update
							oldArb.Update(newArb.Contacts, newArb.NumContacts)
						}

					} else {
						if oldArb != nil {
							//remove the arbiter
							space.Arbiters = append(space.Arbiters[:index], space.Arbiters[index+1:]...)
						}
						newArb.Delete()
					}

				}
			}


		}
	}
}
