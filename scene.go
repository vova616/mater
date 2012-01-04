package mater

import (
	"mater/collision"
)

type Scene struct {
	Space    *collision.Space
	Camera   *Camera
	Entities map[int]*Entity
}

func (scene *Scene) Init(mater *Mater) {
	scene.Space = collision.NewSpace()
	cam := mater.DefaultCamera
	scene.Camera = &cam
	scene.Entities = make(map[int]*Entity, 32)
}

func (scene *Scene) Update(dt float64) {
	scene.Space.Step(dt)
	for _, entity := range scene.Entities {
		if entity.Enabled {
			entity.Update(dt)
		}
	}
}

func (scene *Scene) AddEntity(entity *Entity) {
	scene.Entities[entity.Id()] = entity
}

func (scene *Scene) RemoveEntity(entity *Entity) {
	delete(scene.Entities, entity.Id())
}

func (scene *Scene) Destroy() {
	for _, entity := range scene.Entities {
		entity.Destroy()
	}
}
