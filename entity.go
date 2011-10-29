package mater

import (
	"box2d"
	"github.com/abneptis/GoUUID"
)

type Entity struct {
	id *uuid.UUID
	EntityClass
	Body *box2d.Body
	Enabled bool
	Scene *Scene
}

func (entity *Entity) Init (entityClass EntityClass, scene *Scene) {
	if entity.id != nil {
		return
	}

	entity.id = uuid.NewV4()
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

func (entity *Entity) Id () uuid.UUID {
	return *entity.id
}

type EntityClass interface {
	Init (entity *Entity)
	Update (dt float64)
	Destroy()
}