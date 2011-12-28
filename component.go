package mater

import (
	"os"
	"reflect"
)

type Component interface {
	//used to identify the component
	Name () string
	//Called when added to an entitity
	Init (owner *Entity)
	//Called once per frame if owner.Enabled is true
	Update (owner *Entity, dt float64)
	//Called when removed from an Entity or before the entity is destroyed
	Destroy (owner *Entity)

	//
	Marshal(owner *Entity) ([]byte, os.Error)

	Unmarshal(owner *Entity, data []byte) (os.Error)
}

func (entity *Entity) AddComponent(component Component) {
	name := component.Name()
	//destroy the old component if there is one
	if c2, ok := entity.Components[name]; ok {
		c2.Destroy(entity)
	}
	entity.Components[name] = component
	component.Init(entity)
}

func (entity *Entity) RemoveComponent(component Component) {
	component.Destroy(entity)
	entity.Components[component.Name()] = nil, false
}

func (entity *Entity) RemoveComponentName(name string) {
	if component, ok := entity.Components[name]; ok {
		component.Destroy(entity)
	}
	entity.Components[name] = nil, false
}

var components = make(map[string]reflect.Type)

func RegisterComponent(component Component) {
	components[component.Name()] = reflect.Indirect(reflect.ValueOf(component)).Type()
}

func NewComponent(name string) Component {
	compType, ok := components[name]
	if ok == false {
		dbg.Printf("Error loading component \"%v\", not registered!", name)
		return nil
	}

	component := reflect.New(compType)
	return component.Interface().(Component)
}
