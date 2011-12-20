package mater

import (
	"os"
	"reflect"
	"fmt"
)

type Component interface {
	//used to identify the component
	Name () string
	//Only called when creating a new component at runtime
	//Unmarshalled Components have to call it themselves
	Init (owner *Entity)
	//Called once per frame if owner.Enabled is true
	Update (owner *Entity, dt float64)
	//Called when removed from an Entity or before the entity is destroyed
	Destroy (owner *Entity)

	//
	MarshalJSON(owner *Entity) ([]byte, os.Error)

	UnmarshalJSON(owner *Entity, data []byte) (os.Error)
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

var components = make(map[string]reflect.Type)

func RegisterComponent(component Component) {
	components[component.Name()] = reflect.Indirect(reflect.ValueOf(component)).Type()
}

func NewComponent(name string) Component {
	compType, ok := components[name]
	if ok == false {
		return nil
	}

	component := reflect.New(compType)
	return component.Interface().(Component)
}
