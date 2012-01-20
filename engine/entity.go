package engine

import (
	"github.com/teomat/mater/transform"
	"log"
)

type Entity struct {
	id            int
	Enabled       bool
	Scene         *Scene `json:"-,omitempty"`
	Components    map[string]Component
	ComponentList []Component
	Transform     *transform.Transform
}

var lastEntityId = 0

func nextId() int {
	lastEntityId++
	return lastEntityId
}

func NewEntity() *Entity {
	entity := new(Entity)
	entity.Init()
	return entity
}

func (entity *Entity) Init() {
	if entity.id > 0 {
		log.Printf("Error: Entity already initialized")
		return
	}
	entity.id = nextId()
	entity.Components = make(map[string]Component, 4)
	entity.ComponentList = make([]Component, 0, 4)
}

func (entity *Entity) Update(dt float64) {
	for _, component := range entity.ComponentList {
		component.Update(entity, dt)
	}
}

func (entity *Entity) Destroy() {
	entity.Enabled = false

	for _, component := range entity.Components {
		component.Destroy(entity)
	}

	entity.Scene.RemoveEntity(entity)
	entity.Scene = nil
}

func (entity *Entity) Id() int {
	return entity.id
}

func (entity *Entity) AddComponent(component Component) {
	name := component.Name()
	//destroy the old component if there is one
	if c2, ok := entity.Components[name]; ok {
		entity.RemoveComponent(c2)
	}

	entity.Components[name] = component
	entity.ComponentList = append(entity.ComponentList, component)
	component.Init(entity)

	for _, comp := range entity.ComponentList {
		if comp == component {
			continue
		}

		comp.OnNewComponent(entity, component)
	}

	onNewComponent := entity.Scene.Callbacks.OnNewComponent
	if onNewComponent != nil {
		onNewComponent(entity, component)
	}
}

func (entity *Entity) RemoveComponent(component Component) {
	component.Destroy(entity)

	delete(entity.Components, component.Name())

	i := 0
	cl := entity.ComponentList
	for i = 0; i < len(cl); i++ {
		if component == cl[i] {
			break
		}
	}

	entity.ComponentList = append(cl[:i], cl[i+1:]...)
}

func (entity *Entity) RemoveComponentName(name string) {
	if component, ok := entity.Components[name]; ok {
		entity.RemoveComponent(component)
	}
}
