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
