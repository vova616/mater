package collision

import (
	"mater/vect"
	"math"
	"log"
)

func collide(contacts *[max_points]Contact, sA, sB *Shape) int {
	switch sA.ShapeType() {
		case ShapeType_Circle:
			switch sB.ShapeType() {
				case ShapeType_Circle:
					return collideCircles(contacts, sA, sB, sA.ShapeClass.(*CircleShape), sB.ShapeClass.(*CircleShape))
				default:
					log.Printf("Warning: ShapeB unknown shapetype")
					return 0
			}
		default:
			log.Printf("Warning: ShapeA unknown shapetype")
			return 0
	}
	return 0
}

func collideCircles(contacts *[max_points]Contact, sA, sB *Shape, csA, csB *CircleShape) int {

	xfA := sA.Body.Transform
	xfB := sB.Body.Transform

	p1 := xfA.TransformVect(csA.Position)
	p2 := xfB.TransformVect(csB.Position)

	return circle2circleQuery(p1, p2, csA.Radius, csB.Radius, &contacts[0])
}

func circle2circleQuery(p1, p2 vect.Vect, r1, r2 float64, con *Contact) int {
	minDist := r1 + r2

	delta := vect.Sub(p2, p1)
	distSqr := delta.LengthSqr()

	if distSqr >= minDist * minDist {
		return 0
	}

	dist := math.Sqrt(distSqr)

	con.Separation = dist - minDist
	pDist := dist
	if dist == 0.0 {
		pDist = math.Inf(1)
	}

	pos := vect.Add(p1, vect.Mult(delta, 0.5 + (r1 - 0.5 * minDist)/pDist))

	norm := vect.Vect{1, 0}

	if dist != 0.0 {
		norm = vect.Mult(delta, 1.0 / dist)
	}

	norm.Mult(-1)

	con.Reset(pos, norm, dist - minDist)

	con.R1 = vect.Sub(con.Position, p1)
	con.R2 = vect.Sub(con.Position, p2)

	return 1
}

func segmentEncapQuery(p1, p2 vect.Vect, r1, r2 float64, con *Contact, tangent vect.Vect) int {
	count := circle2circleQuery(p1, p2, r1, r2, con)
	if vect.Dot(con.Normal, tangent) >= 0.0 {
		return count
	} else {
		return 0
	}
	panic("Never reached")
}