package collision

import (
	"mater/vect"
	"log"
)

type World struct {
	Gravity vect.Vect
	StaticBodies []*Body
	DynamicBodies []*Body
}

func (world *World) init () {
	world.StaticBodies = make([]*Body, 0, 8)
	world.DynamicBodies = make([]*Body, 0, 8)
}

func NewWorld() *World {
	world := new(World)
	world.init()

	return world
}

func (world *World) AddStaticBody (body *Body) {
	if body.World != nil {
		log.Printf("Error adding static body: body.World != nil")
		return
	}
	body.World = world
	world.StaticBodies = append(world.StaticBodies, body)
}

func (world *World) AddDynamicBody (body *Body) {
	if body.World != nil {
		log.Printf("Error adding dynamic body: body.World != nil")
		return
	}
	body.World = world
	world.DynamicBodies = append(world.DynamicBodies, body)
}

func (world *World) RemoveStaticBody (body *Body) {
	bodies := world.StaticBodies
	for i, b := range bodies {
		if b == body {
			world.StaticBodies = append(bodies[:i], bodies[i+1:]...)
			return
		}
	}
	log.Printf("Warning removing body: static body not found!")
}

func (world *World) RemoveDynamicBody (body *Body) {
	bodies := world.StaticBodies
	for i, b := range bodies {
		if b == body {
			world.StaticBodies = append(bodies[:i], bodies[i+1:]...)
			return
		}
	}
	log.Printf("Warning removing body: dynamic body not found!")
}
