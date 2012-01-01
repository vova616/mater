package collision

import (
	"mater/vect"
	"log"
)

type Space struct {
	Enabled bool
	Gravity vect.Vect
	StaticBodies []*Body
	DynamicBodies []*Body
}

func (space *Space) init() {
	space.StaticBodies = make([]*Body, 0, 8)
	space.DynamicBodies = make([]*Body, 0, 8)
	space.Enabled = true
}

func NewSpace() *Space {
	space := new(Space)
	space.init()
	
	return space
}

func (space *Space) AddBody(body *Body) {
	if body.IsStatic() {
		if body.Space != nil {
			log.Printf("Error adding static body: body.Space != nil")
			return
		}
		body.Space = space
		space.StaticBodies = append(space.StaticBodies, body)
	} else {
		if body.Space != nil {
			log.Printf("Error adding dynamic body: body.Space != nil")
			return
		}
		body.Space = space
		space.DynamicBodies = append(space.DynamicBodies, body)
	}
}

func (space *Space) RemoveBody(body *Body) {
	if body.IsStatic() {
		bodies := space.StaticBodies
		for i, b := range bodies {
			if b == body {
				space.StaticBodies = append(bodies[:i], bodies[i+1:]...)
				return
			}
		}
		log.Printf("Warning removing body: static body not found!")
	} else {
		bodies := space.StaticBodies
		for i, b := range bodies {
			if b == body {
				space.StaticBodies = append(bodies[:i], bodies[i+1:]...)
				return
			}
		}
		log.Printf("Warning removing body: dynamic body not found!")
	}
}

func (space *Space) Step(dt float64) {
	
}
