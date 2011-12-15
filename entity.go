package mater

import (
	"box2d"
)

type Entity struct {
	id int
	Body *box2d.Body `json:",omitempty"`
	Enabled bool
	Scene *Scene `json:"-,omitempty"`
	Components map[string]Component
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

	entity.Body = nil
	entity.Scene.RemoveEntity(entity)
	entity.Scene = nil
}

func (entity *Entity) Id () int {
	return entity.id
}
