package mater

import (
	"box2d"
)

type Scene struct {
	World *box2d.World
	Camera *Camera
	Entities map[int]*Entity
}

func (scene *Scene) Init (mater *Mater) {
	scene.World = new(box2d.World)
	cam := mater.DefaultCamera
	scene.Camera = &cam
	scene.World.Init()
	scene.Entities = make(map[int]*Entity, 32)
}

func (scene *Scene) Update (dt float64) {
	scene.World.Step(dt)
}

func (scene *Scene) AddEntity(entity *Entity) {
	scene.Entities[entity.Id()] = entity
}

func (scene *Scene) RemoveEntity(entity *Entity) {
	scene.Entities[entity.Id()] = nil, false
}
