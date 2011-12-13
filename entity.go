package mater

import (
	"box2d"
)

type Entity struct {
	id int
	Body *box2d.Body `json:",omitempty"`
	Enabled bool
	Scene *Scene `json:"-,omitempty"`
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
}

func (entity *Entity) Destroy () {
	entity.Scene = nil
	entity.Enabled = false
}

func (entity *Entity) Id () int {
	return entity.id
}

type EntityClass interface {
	Init (entity *Entity)
	Update (dt float64)
	Destroy()
}
