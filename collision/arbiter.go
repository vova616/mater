// Package mater/collision procies a basic 2d collision library,
// based on box2d-lite and chipmunk-physics.

package collision

import (
	"github.com/teomat/mater/vect"
	"math"
)

// The maximum number of ContactPoints a single Arbiter can have.
const MaxPoints = 2

type Arbiter struct {
	ShapeA, ShapeB *Shape
	Contacts       [MaxPoints]Contact
	NumContacts    int

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

	arb.NumContacts = collide(&arb.Contacts, arb.ShapeA, arb.ShapeB)

	arb.Friction = math.Sqrt(sa.Friction * sb.Friction)

	return arb
}

func (arb *Arbiter) delete() {
	arb.ShapeA = nil
	arb.ShapeB = nil
	arb.NumContacts = 0
	arb.Friction = 0
}

func (arb1 *Arbiter) equals(arb2 *Arbiter) bool {
	if arb1.ShapeA == arb2.ShapeA && arb1.ShapeB == arb2.ShapeB {
		return true
	}
	/*if arb1.ShapeA == arb2.ShapeB && arb1.ShapeB == arb1.ShapeA {
		return true
	}*/

	return false
}

/*
func (arb *Arbiter) Update(newContacts [MaxPoints]Contact, numNewContacts int) {
	var mergedContacts [MaxPoints]Contact

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
}*/

func (arb *Arbiter) preStep(inv_dt float64) {
	const allowedPenetration = 0.01
	biasFactor := 0.0
	if Settings.PositionCorrection {
		biasFactor = 0.2
	}

	b1 := arb.ShapeA.Body
	b2 := arb.ShapeB.Body

	for i := 0; i < arb.NumContacts; i++ {
		c := &arb.Contacts[i]

		c.R1 = vect.Sub(c.Position, b1.Transform.Position)
		c.R2 = vect.Sub(c.Position, b2.Transform.Position)
		r1 := c.R1
		r2 := c.R2

		//Precompute normal mass, tangent mass, and bias
		rn1 := vect.Dot(r1, c.Normal)
		rn2 := vect.Dot(r2, c.Normal)
		kNormal := b1.invMass + b2.invMass
		kNormal += b1.invI*(vect.Dot(r1, r1)-rn1*rn1) +
			b2.invI*(vect.Dot(r2, r2)-rn2*rn2)
		c.MassNormal = 1.0 / kNormal

		tangent := vect.CrossVF(c.Normal, 1.0)
		rt1 := vect.Dot(r1, tangent)
		rt2 := vect.Dot(r2, tangent)
		kTangent := b1.invMass + b2.invMass
		kTangent += b1.invI*(vect.Dot(r1, r1)-rt1*rt1) +
			b2.invI*(vect.Dot(r2, r2)-rt2*rt2)
		c.MassTangent = 1.0 / kTangent

		c.Bias = -biasFactor * inv_dt * math.Min(0.0, c.Separation+allowedPenetration)

		if Settings.AccumulateImpulses {
			//Apply normal + friction impulse
			P := vect.Add(vect.Mult(c.Normal, c.Pn), vect.Mult(tangent, c.Pt))

			b1.Velocity.Sub(vect.Mult(P, b1.invMass))
			b1.AngularVelocity -= b1.invI * vect.Cross(r1, P)

			b2.Velocity.Add(vect.Mult(P, b2.invMass))
			b2.AngularVelocity += b2.invI * vect.Cross(r2, P)
		}

	}
}

func (arb *Arbiter) applyImpulse() {
	sA := arb.ShapeA
	sB := arb.ShapeB

	b1 := sA.Body
	b2 := sB.Body

	//xfA := b1.Transform
	//xfB := b2.Transform

	for i := 0; i < arb.NumContacts; i++ {
		c := &arb.Contacts[i]

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
			c.Pn = math.Max(Pn0+dPn, 0.0)
			dPn = c.Pn - Pn0
		} else {
			dPn = math.Max(dPn, 0.0)
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
			c.Pt = clamp(oldTangentImpulse+dPt, -maxPt, maxPt)
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
