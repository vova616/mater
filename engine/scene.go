package engine

import (
	"github.com/teomat/mater/collision"
)

type Callbacks struct {
	OnNewComponent func(entity *Entity, comp Component)
}

type Scene struct {
	Space        *collision.Space
	StaticEntity Entity
	Entities     map[int]*Entity
	Callbacks    Callbacks
}

func (scene *Scene) Init() {
	scene.Space = collision.NewSpace()
	scene.StaticEntity.Init()
	scene.StaticEntity.Scene = scene
	scene.Entities = make(map[int]*Entity, 32)
}

func (scene *Scene) Update(dt float64) {
	scene.Space.Step(dt)
	scene.StaticEntity.Update(dt)
	for _, entity := range scene.Entities {
		if entity.Enabled {
			entity.Update(dt)
		}
	}
}

func (scene *Scene) AddEntity(entity *Entity) {
	entity.Scene = scene
	scene.Entities[entity.Id()] = entity
}

func (scene *Scene) RemoveEntity(entity *Entity) {
	delete(scene.Entities, entity.Id())
	entity.Destroy()
}

func (scene *Scene) Destroy() {
	for _, entity := range scene.Entities {
		entity.Destroy()
	}
}
