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

var lastId = 0
func nextId() int {
	lastId++
	return lastId
}

func (entity *Entity) Init (scene *Scene) {
	if entity.id != 0 {
		return
	}

	entity.id = nextId()
	entity.Enabled = true
	entity.Scene = scene
	entity.Components = make(map[string]Component)
}

func (entity *Entity) Update (dt float64) {
	for _, component := range entity.Components {
		component.Update(entity, dt)
	}
}

func (entity *Entity) Destroy () {
	entity.Scene = nil
	entity.Enabled = false

	for _, component := range entity.Components {
		component.Destroy(entity)
	}
}

func (entity *Entity) Id () int {
	return entity.id
}
