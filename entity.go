package mater

import (
	"github.com/teomat/mater/transform"
)

type Entity struct {
	id int
	Enabled bool
	Scene *Scene `json:"-,omitempty"`
	Components map[string]Component
	Transform *transform.Transform
}

var lastEntityId = 0
func nextId() int {
	lastEntityId++
	return lastEntityId
}

func NewEntity () *Entity {
	entity := new(Entity)
	entity.id = nextId()
	entity.Components = make(map[string]Component)
	return entity
}

func (entity *Entity) Update (dt float64) {
	for _, component := range entity.Components {
		component.Update(entity, dt)
	}
}

func (entity *Entity) Destroy () {
	entity.Enabled = false

	for _, component := range entity.Components {
		component.Destroy(entity)
	}

	entity.Scene.RemoveEntity(entity)
	entity.Scene = nil
}

func (entity *Entity) Id () int {
	return entity.id
}
