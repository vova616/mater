package collision

import (
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
	"math"
)

// Used to keep a linked list of all arbiters on a body.
type ArbiterEdge struct {
	Arbiter    *Arbiter
	Next, Prev *ArbiterEdge
	Other      *Body
}

type arbiterState int

const (
	arbiterStateFirstColl = iota
	arbiterStateNormal
)

// The maximum number of ContactPoints a single Arbiter can have.
const MaxPoints = 2

type Arbiter struct {
	// The two colliding shapes.
	ShapeA, ShapeB *Shape
	// The contact points between the shapes.
	Contacts [MaxPoints]Contact
	// The number of contact points.
	NumContacts int

	nodeA, nodeB *ArbiterEdge

	Friction    float64
	Restitution float64

	Surface_vr vect.Vect

	// Used to keep a linked list of all arbiters in a space.
	Next, Prev *Arbiter

	state arbiterState
}

func newArbiter() *Arbiter {
	return new(Arbiter)
}

// Creates an arbiter between the given shapes.
// If the shapes do not collide, arbiter.NumContact is zero.
func CreateArbiter(sa, sb *Shape) *Arbiter {
	arb := newArbiter()

	if sa.ShapeType() < sb.ShapeType() {
		arb.ShapeA = sa
		arb.ShapeB = sb
	} else {
		arb.ShapeA = sb
		arb.ShapeB = sa
	}

	arb.Surface_vr = vect.Vect{}

	arb.nodeA = new(ArbiterEdge)
	arb.nodeB = new(ArbiterEdge)

	return arb
}

func (arb *Arbiter) destroy() {
	arb.ShapeA = nil
	arb.ShapeB = nil
	arb.NumContacts = 0
	arb.Friction = 0
}

func (arb1 *Arbiter) equals(arb2 *Arbiter) bool {
	if arb1.ShapeA == arb2.ShapeA && arb1.ShapeB == arb2.ShapeB {
		return true
	}

	return false
}

func (arb *Arbiter) update(contacts *[MaxPoints]Contact, numContacts int) {
	oldContacts := &arb.Contacts
	oldNumContacts := arb.NumContacts

	sa := arb.ShapeA
	sb := arb.ShapeB

	for i := 0; i < oldNumContacts; i++ {
		oldC := &oldContacts[i]
		for j := 0; j < numContacts; j++ {
			newC := &contacts[j]

			if newC.hash == oldC.hash {
				newC.jnAcc = oldC.jnAcc
				newC.jtAcc = oldC.jtAcc
			}
		}
	}

	arb.Contacts = *contacts
	arb.NumContacts = numContacts

	arb.Friction = sa.Friction * sb.Friction
	arb.Restitution = sa.Restitution * sb.Restitution

	arb.Surface_vr = vect.Sub(sa.Surface_v, sb.Surface_v)
}

func (arb *Arbiter) preStep(inv_dt float64, slop, bias float64) {

	a := arb.ShapeA.Body
	b := arb.ShapeB.Body

	for i := 0; i < arb.NumContacts; i++ {
		con := &arb.Contacts[i]

		// Calculate the offsets.
		con.R1 = vect.Sub(con.Position, a.Transform.Position)
		con.R2 = vect.Sub(con.Position, b.Transform.Position)

		// Calculate the mass normal and mass tangent.
		con.nMass = 1.0 / k_scalar(a, b, con.R1, con.R2, con.Normal)
		con.tMass = 1.0 / k_scalar(a, b, con.R1, con.R2, vect.Perp(con.Normal))

		// Calculate the target bias velocity.
		con.bias = -bias * inv_dt * math.Min(0.0, con.Dist+slop)
		con.jBias = 0.0

		// Calculate the target bounce velocity.
		con.bounce = normal_relative_velocity(a, b, con.R1, con.R2, con.Normal) * arb.Restitution
	}
}

func (arb *Arbiter) applyCachedImpulse(dt_coef float64) {
	if arb.state == arbiterStateFirstColl && arb.NumContacts > 0 {
		arb.state = arbiterStateNormal
		return
	}

	a := arb.ShapeA.Body
	b := arb.ShapeB.Body
	for i := 0; i < arb.NumContacts; i++ {
		con := &arb.Contacts[i]
		j := transform.RotateVect(con.Normal, transform.Rotation{con.jnAcc, con.jtAcc})
		apply_impulses(a, b, con.R1, con.R2, vect.Mult(j, dt_coef))
	}
}

func (arb *Arbiter) applyImpulse() {
	a := arb.ShapeA.Body
	b := arb.ShapeB.Body

	for i := 0; i < arb.NumContacts; i++ {
		con := &arb.Contacts[i]
		n := con.Normal
		r1 := con.R1
		r2 := con.R2

		// Calculate the relative bias velocities.
		vb1 := vect.Add(a.v_bias, vect.Mult(vect.Perp(r1), a.w_bias))
		vb2 := vect.Add(b.v_bias, vect.Mult(vect.Perp(r2), b.w_bias))
		vbn := vect.Dot(vect.Sub(vb2, vb1), n)

		// Calculate and clamp the bias impulse.
		jbn := (con.bias - vbn) * con.nMass
		jbnOld := con.jBias
		con.jBias = math.Max(jbnOld+jbn, 0.0)
		jbn = con.jBias - jbnOld

		// Apply the bias impulse.
		apply_bias_impulses(a, b, r1, r2, vect.Mult(n, jbn))

		// Calculate the relative velocity.
		vr := relative_velocity(a, b, r1, r2)
		vrn := vect.Dot(vr, n)

		// Calculate and clamp the normal impulse.
		jn := -(con.bounce + vrn) * con.nMass
		jnOld := con.jnAcc
		con.jnAcc = math.Max(jnOld+jn, 0.0)
		jn = con.jnAcc - jnOld

		// Calculate the relative tangent velocity.
		vrt := vect.Dot(vect.Add(vr, arb.Surface_vr), vect.Perp(n))

		// Calculate and clamp the friction impulse.
		jtMax := arb.Friction * con.jnAcc
		jt := -vrt * con.tMass
		jtOld := con.jtAcc
		con.jtAcc = clamp(jtOld+jt, -jtMax, jtMax)
		jt = con.jtAcc - jtOld

		// Apply the final impulse.
		apply_impulses(a, b, r1, r2, transform.RotateVect(n, transform.Rotation{jn, jt}))
	}
}
