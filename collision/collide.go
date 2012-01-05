package collision

import (
	"log"
	"mater/vect"
	"math"
)

type collisionHandler func(contacts *[max_points]Contact, sA, sB *Shape) int

var collisionHandlers = [numShapes][numShapes]collisionHandler{
	ShapeType_Circle: [numShapes]collisionHandler{
		ShapeType_Circle:  circle2circle,
		ShapeType_Segment: circle2segment,
		ShapeType_Polygon: circle2polygon,
	},
	ShapeType_Segment: [numShapes]collisionHandler{
		ShapeType_Circle:  nil,
		ShapeType_Segment: nil,
		ShapeType_Polygon: nil,
	},
	ShapeType_Polygon: [numShapes]collisionHandler{
		ShapeType_Circle: nil,
		ShapeType_Segment: nil,
		ShapeType_Polygon: nil,
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

//START COLLISION HANDLERS
func circle2circle(contacts *[max_points]Contact, sA, sB *Shape) int {
	csA, ok := sA.ShapeClass.(*CircleShape)
	if !ok {
		log.Printf("Error: ShapeA not a CircleShape!")
		return 0
	}
	csB, ok := sB.ShapeClass.(*CircleShape)
	if !ok {
		log.Printf("Error: ShapeA not a CircleShape!")
		return 0
	}
	return circle2circleQuery(csA.Tc, csB.Tc, csA.Radius, csB.Radius, &contacts[0])
}

func circle2segment(contacts *[max_points]Contact, sA, sB *Shape) int {
	circle, ok := sA.ShapeClass.(*CircleShape)
	if !ok {
		log.Printf("Error: ShapeA not a CircleShape!")
		return 0
	}
	segment, ok := sB.ShapeClass.(*SegmentShape)
	if !ok {
		log.Printf("Error: ShapeB not a SegmentShape!")
		return 0
	}

	return circle2segmentFunc(contacts, circle, segment)
}

func circle2polygon(contacts *[max_points]Contact, sA, sB *Shape) int {
	circle, ok := sA.ShapeClass.(*CircleShape)
	if !ok {
		log.Printf("Error: ShapeA not a CircleShape!")
		return 0
	}
	poly, ok := sB.ShapeClass.(*PolygonShape)
	if !ok {
		log.Printf("Error: ShapeB not a SegmentShape!")
		return 0
	}

	return circle2polyFunc(contacts, circle, poly)
}
//END COLLISION HANDLERS

func circle2circleQuery(p1, p2 vect.Vect, r1, r2 float64, con *Contact) int {
	minDist := r1 + r2

	delta := vect.Sub(p2, p1)
	distSqr := delta.LengthSqr()

	if distSqr >= minDist*minDist {
		return 0
	}

	dist := math.Sqrt(distSqr)

	con.Separation = dist - minDist
	pDist := dist
	if dist == 0.0 {
		pDist = math.Inf(1)
	}

	pos := vect.Add(p1, vect.Mult(delta, 0.5+(r1-0.5*minDist)/pDist))

	norm := vect.Vect{1, 0}

	if dist != 0.0 {
		norm = vect.Mult(delta, 1.0/dist)
	}

	con.Reset(pos, norm, dist-minDist)

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

func circle2segmentFunc(contacts *[max_points]Contact, circle *CircleShape, segment *SegmentShape) int {
	rsum := circle.Radius + segment.Radius

	//Calculate normal distance from segment
	dn := vect.Dot(segment.Tn, circle.Tc) - vect.Dot(segment.Ta, segment.Tn)
	dist := math.Abs(dn) - rsum
	if dist > 0.0 {
		return 0
	}

	//Calculate tangential distance along segment
	dt := -vect.Cross(segment.Tn, circle.Tc)
	dtMin := -vect.Cross(segment.Tn, segment.Ta)
	dtMax := -vect.Cross(segment.Tn, segment.Tb)

	// Decision tree to decide which feature of the segment to collide with.
	if dt < dtMin {
		if dt < (dtMin - rsum) {
			return 0
		} else {
			return segmentEncapQuery(circle.Tc, segment.Ta, circle.Radius, segment.Radius, &contacts[0], segment.A_tangent)
		}
	} else {
		if dt < dtMax {
			n := segment.Tn
			if dn >= 0.0 {
				n.Mult(-1)
			}
			con := &contacts[0]
			pos := vect.Add(circle.Tc, vect.Mult(n, circle.Radius+dist*0.5))
			con.Reset(pos, n, dist)
			return 1
		} else {
			if dt < (dtMax + rsum) {
				return segmentEncapQuery(circle.Tc, segment.Tb, circle.Radius, segment.Radius, &contacts[0], segment.B_tangent)
			} else {
				return 0
			}
		}
	}
	panic("Never reached")
}

func circle2polyFunc(contacts *[max_points]Contact, circle *CircleShape, poly *PolygonShape) int {
	
	axes := poly.TAxes

	mini := 0
	min := vect.Dot(axes[0].N, circle.Tc) - axes[0].D - circle.Radius
	for i, axis := range axes {
		dist := vect.Dot(axis.N, circle.Tc) - axis.D - circle.Radius
		if dist > 0.0 {
			return 0
		} else if dist > min {
			min = dist
			mini = i
		}
	}

	n := axes[mini].N
	a := poly.TVerts[mini]
	b := poly.TVerts[(mini + 1) % poly.NumVerts]
	dta := vect.Cross(n, a)
	dtb := vect.Cross(n, b)
	dt := vect.Cross(n, circle.Tc)

	if dt < dtb {
		return circle2circleQuery(circle.Tc, b, circle.Radius, 0.0, &contacts[0])
	} else if dt < dta {
		contacts[0].Reset(
			vect.Sub(circle.Tc, vect.Mult(n, circle.Radius + min / 2.0)),
			vect.Mult(n, -1),
			min,
		)
		return 1
	} else {
		return circle2circleQuery(circle.Tc, a, circle.Radius, 0.0, &contacts[0])
	}
	panic("Never reached")
}
