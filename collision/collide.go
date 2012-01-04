package collision

import (
	"mater/vect"
	"math"
	"log"
)

type collisionHandler func(contacts *[max_points]Contact, sA, sB *Shape) int
var collisionHandlers = [numShapes][numShapes]collisionHandler{
	ShapeType_Circle: [numShapes]collisionHandler{
		ShapeType_Circle: circle2circle,
		ShapeType_Segment: circle2segment,
	},
	ShapeType_Segment: [numShapes]collisionHandler{
		ShapeType_Circle: nil,
		ShapeType_Segment: nil,
	},
}

func collide(contacts *[max_points]Contact, sA, sB *Shape) int {
	stA := sA.ShapeType()
	stB := sB.ShapeType()

	if stA > stB {
		log.Printf("Error: shapes not ordered")
		return 0
	}

	handler := collisionHandlers[stA][stB]
	if handler == nil {
		return 0
	}

	return handler(contacts, sA, sB)
}

func circle2circle(contacts *[max_points]Contact, sA, sB *Shape) int {
	csA, ok := sA.ShapeClass.(*CircleShape)
	if ! ok {
		log.Printf("Error: ShapeA not a CircleShape!")
		return 0
	}
	csB, ok := sB.ShapeClass.(*CircleShape)
	if ! ok {
		log.Printf("Error: ShapeA not a CircleShape!")
		return 0
	}
	return circle2circleQuery(csA.tc, csB.tc, csA.Radius, csB.Radius, &contacts[0])
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

//circle-segment collision taken from chipmunk-physics
func circle2segment(contacts *[max_points]Contact, sA, sB *Shape) int {
	circle, ok := sA.ShapeClass.(*CircleShape)
	if ! ok {
		log.Printf("Error: ShapeA not a CircleShape!")
		return 0
	}
	segment, ok := sB.ShapeClass.(*SegmentShape)
	if ! ok {
		log.Printf("Error: ShapeB not a SegmentShape!")
		return 0
	}

	rsum := circle.Radius + segment.Radius

	//Calculate normal distance from segment
	dn := vect.Dot(segment.tn, circle.tc) - vect.Dot(segment.ta, segment.tn)
	dist := math.Fabs(dn) - rsum
	if dist > 0.0 {
		return 0
	}

	//Calculate tangential distance along segment
	dt := -vect.Cross(segment.tn, circle.tc)
	dtMin := -vect.Cross(segment.tn, segment.ta)
	dtMax := -vect.Cross(segment.tn, segment.tb)

	// Decision tree to decide which feature of the segment to collide with.
	if dt < dtMin {
		if dt < (dtMin - rsum) {
			return 0
		} else {
			return segmentEncapQuery(circle.tc, segment.ta, circle.Radius, segment.Radius, &contacts[0], segment.a_tangent)
		}
	} else {
		if dt < dtMax {
			n := segment.tn
			if dn >= 0.0 {
				n.Mult(-1)
			}
			con := &contacts[0]
			pos := vect.Add(circle.tc, vect.Mult(n, circle.Radius + dist * 0.5))
			con.Reset(pos, n, dist)
			return 1
		} else {
			if dt < (dtMax + rsum) {
				return segmentEncapQuery(circle.tc, segment.tb, circle.Radius, segment.Radius, &contacts[0], segment.b_tangent)
			} else {
				return 0
			}
		}
	}
	panic("Never reached")
}
