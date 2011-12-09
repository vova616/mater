package mater

import (
	"box2d"
)

type Entity struct {
	id int
	EntityClass
	Body *box2d.Body
	Enabled bool
	Scene *Scene
}

var lastId = 0
func nextId() int {
	lastId++
	return lastId
}

func (entity *Entity) Init (entityClass EntityClass, scene *Scene) {
	if entity.id != 0 {
		return
	}

	entity.id = nextId()
	entity.Enabled = true
	entity.Scene = scene

	entity.EntityClass = entityClass
	entityClass.Init(entity)
}

func (entity *Entity) Destroy () {
	entity.EntityClass.Destroy()

	entity.Scene = nil
	entity.EntityClass = nil
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