package collision

import (
	"math"
	"mater/vect"
)

const max_points = 2

type Arbiter struct {
	ShapeA, ShapeB *Shape
	Contacts [max_points]Contact
	NumContacts int

	Friction float64
}

func newArbiter() *Arbiter {
	return new(Arbiter)
}

func CreateArbiter(sa, sb *Shape) *Arbiter {
	arb := newArbiter()

	if sa.ShapeType() < sb.ShapeType() {
		arb.ShapeA = sa
		arb.ShapeB = sb
	} else {
		arb.ShapeA = sb
		arb.ShapeB = sa
	}

	arb.NumContacts = collide(&arb.Contacts, sa, sb)

	arb.Friction = math.Sqrt(sa.Friction * sb.Friction)

	return arb
}

func (arb *Arbiter) Delete () {
	
}

func (arb1 *Arbiter) Equals(arb2 *Arbiter) bool {
	if arb1.ShapeA == arb2.ShapeA && arb1.ShapeB == arb2.ShapeB {
		return true
	}
	if arb1.ShapeA == arb2.ShapeB && arb1.ShapeB == arb1.ShapeA {
		return true
	}

	return false
}

func (arb *Arbiter) Update(newContacts [max_points]Contact, numNewContacts int) {
	var mergedContacts [max_points]Contact

	for i := 0; i < numNewContacts; i++ {
		cNew := newContacts[i]
		k := -1
		for j := 0; j < arb.NumContacts; j++ {
			cOld := arb.Contacts[j]

			if cNew.Feature.Value() == cOld.Feature.Value() {
				k = j
				break
			}
		}

		if k > -1 {
			cOld := arb.Contacts[k]
			mergedContacts[i] = cNew
			c := mergedContacts[i]
			const warmStarting = false
			if warmStarting {
				c.Pn = cOld.Pn
				c.Pt = cOld.Pt
				c.Pnb = cOld.Pnb
			} else {
				c.Pn = 0
				c.Pt = 0
				c.Pnb = 0
			}
		} else {
			mergedContacts[i] = newContacts[i]
		}
	}

	for i := 0; i < numNewContacts; i++ {
		arb.Contacts[i] = mergedContacts[i]
	}

	arb.NumContacts = numNewContacts
}

func (arb *Arbiter) ApplyImpulse() {
	sA := arb.ShapeA
	sB := arb.ShapeB

	b1 := sA.Body
	b2 := sB.Body

	//xfA := b1.Transform
	//xfB := b2.Transform

	for i := 0; i < arb.NumContacts; i++ {
		c := arb.Contacts[i]

		// Relative velocity at contact
		dv := vect.Vect{}
		{
			t1 := vect.Add(b2.Velocity, vect.CrossFV(b2.AngularVelocity, c.R2))
			t2 := vect.Sub(b1.Velocity, vect.CrossFV(b1.AngularVelocity, c.R1))

			dv = vect.Sub(t1, t2)
		}

		// Compute normal impulse
		vn := vect.Dot(dv, c.Normal)

		dPn := c.MassNormal * (-vn + c.Bias)

		if Settings.AccumulateImpulses {
			// Clamp the accumulated impulse
			Pn0 := c.Pn
			c.Pn = math.Fmax(Pn0 + dPn, 0.0)
			dPn = c.Pn - Pn0
		} else {
			dPn = math.Fmax(dPn, 0.0)
		}

		//Apply contact impulse
		Pn := vect.Mult(c.Normal, dPn)

		b1.Velocity.Sub(vect.Mult(Pn, b1.invMass))
		b1.AngularVelocity -= b1.invI * vect.Cross(c.R1, Pn)
		
		b2.Velocity.Add(vect.Mult(Pn, b2.invMass))
		b2.AngularVelocity += b2.invI * vect.Cross(c.R2, Pn)
		
		//Relative velocity at contact
		{
			t1 := vect.Add(b2.Velocity, vect.CrossFV(b2.AngularVelocity, c.R2))
			t2 := vect.Sub(b1.Velocity, vect.CrossFV(b1.AngularVelocity, c.R1))

			dv = vect.Sub(t1, t2)
		}

		tangent := vect.CrossVF(c.Normal, 1.0)
		vt := vect.Dot(dv, tangent)
		dPt := c.MassTangent * (-vt)

		if Settings.AccumulateImpulses {
			//Compute friction impulse
			maxPt := arb.Friction * c.Pn

			//Clamp Friction
			oldTangentImpulse := c.Pt
			c.Pt = clamp(oldTangentImpulse + dPt, -maxPt, maxPt)
			dPt = c.Pt - oldTangentImpulse
		} else {
			maxPt := arb.Friction * dPn
			dPt = clamp(dPt, -maxPt, maxPt)
		}

		// Apply contact impulse
		Pt := vect.Mult(tangent, dPt)

		b1.Velocity.Sub(vect.Mult(Pt, b1.invMass))
		b1.AngularVelocity -= b1.invI * vect.Cross(c.R1, Pt)

		b2.Velocity.Add(vect.Mult(Pt, b2.invMass))
		b2.AngularVelocity += b2.invI * vect.Cross(c.R2, Pt)
	}
}

