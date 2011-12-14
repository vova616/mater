package mater

import (
	"os"
)

type Component interface {
	Name () string
	Init (owner *Entity)
	Update (owner *Entity, dt float64)
	Destroy (owner *Entity)
}

func (entity *Entity) AddComponent(component Component) {
	entity.Components[component.Name()] = component
	component.Init(entity)
}

func (entity *Entity) RemoveComponent(component Component) {
	name := component.Name()
	if _, ok := entity.Components[name]; ok {
		component.Destroy(entity)
		entity.Components[name] = nil, false
	}
}

func (entity *Entity) RemoveComponentName(name string) {
	component, ok := entity.Components[name]
	if ok {
		component.Destroy(entity)
		entity.Components[name] = nil, false
	}
}

//For a component to be un/marshalled it has to be registered as a serializable component
type SerializableComponent interface {
	MarshalJSON(owner *Entity) ([]byte, os.Error)
	UnmarshalJSON(owner *Entity, data []byte) os.Error
}

var serializableComponents = make(map[string]SerializableComponent)

func RegisterSerializableComponent(name string, component SerializableComponent) {
	
	serializableComponents[name] = component
}
