package collision

import (
	"mater/vect"
	"log"
)

type Space struct {
	Enabled bool
	Gravity vect.Vect
	Bodies []*Body
}

func (space *Space) init() {
	space.Bodies = make([]*Body, 0, 16)
	space.Enabled = true
}

func NewSpace() *Space {
	space := new(Space)
	space.init()
	
	return space
}

func (space *Space) AddBody(body *Body) {
	if body.Space != nil {
		log.Printf("Error adding body: body.Space != nil")
		return
	}
	body.Space = space
	space.Bodies = append(space.Bodies, body)
}

func (space *Space) RemoveBody(body *Body) {
	bodies := space.Bodies
	for i, b := range bodies {
		if b == body {
			space.Bodies = append(bodies[:i], bodies[i+1:]...)
			return
		}
	}
	log.Printf("Warning removing body: body not found!")
}

func (space *Space) Step(dt float64) {
	
}
