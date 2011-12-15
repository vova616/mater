package mater

import (
	"os"
)

type Component interface {
	//name can be different than the name components register themselves with an entity, but it is then used when unmarshalling them because interfaces cannot be unmarhsalled from json
	Name () string
	Init (owner *Entity)
	Update (owner *Entity, dt float64)
	Destroy (owner *Entity)
}

func (entity *Entity) AddComponent(component Component) {
	//components have to add/remove themselves from entity.Components
	component.Init(entity)
}

func (entity *Entity) RemoveComponent(component Component) {
	//components have to add/remove themselves from entity.Components
	component.Destroy(entity)
}

func (entity *Entity) RemoveComponentName(name string) {
	if component, ok := entity.Components[name]; ok {
		component.Destroy(entity)
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

func UnregisterSerializableComponent(name string) {
	serializableComponents[name] = nil, false
}
