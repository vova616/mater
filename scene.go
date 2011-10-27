package mater

import (
	"box2d"
)

type Scene struct {
	World *box2d.World
	//Entities []*Entity
}

func (scene *Scene) Init () {
	scene.World = new(box2d.World)
	scene.World.Init()
	//scene.Entities = make([]*Entity, 0, 16)
}

func (scene *Scene) Update (dt float64) {
	scene.World.Step(dt)
}