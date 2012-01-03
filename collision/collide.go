package collision

import (
//	"mater/vect"
//	"mater/transform"
	"log"
)

func collide(contacts [max_points]Contact, sA, sB *Shape) int {
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

func collideCircles(contacts [max_points]Contact, sA, sB *Shape, csA, csB *CircleShape) int {
	//mindist := csA.Radius + csB.Radius
	return 0
}