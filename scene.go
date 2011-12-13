package mater

import (
	"box2d"
)

type Scene struct {
	World *box2d.World
	Camera *Camera
	Entities []*Entity
}

func (scene *Scene) Init (mater *Mater) {
	scene.World = new(box2d.World)
	cam := mater.DefaultCamera
	scene.Camera = &cam
	scene.World.Init()
}

func (scene *Scene) Update (dt float64) {
	scene.World.Step(dt)
}
