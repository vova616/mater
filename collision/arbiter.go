package collision

import (
	"math"
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

	arb.NumContacts = collide(arb.Contacts, sa, sb)

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
