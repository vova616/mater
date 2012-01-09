package mater

import (
	"log"
	"reflect"
)

type Component interface {
	//used to identify the component
	Name() string
	//Called when added to an entitity
	Init(owner *Entity)
	//Called once per frame if owner.Enabled is true
	Update(owner *Entity, dt float64)
	//Called when removed from an Entity or before the entity is destroyed
	Destroy(owner *Entity)
	//Called when a new component is added to the owner
	OnNewComponent(owner *Entity, other Component)

	//
	Marshal(owner *Entity) ([]byte, error)

	Unmarshal(owner *Entity, data []byte) error
}

func (entity *Entity) AddComponent(component Component) {
	name := component.Name()
	//destroy the old component if there is one
	if c2, ok := entity.Components[name]; ok {
		c2.Destroy(entity)
	}
	entity.Components[name] = component
	component.Init(entity)
	for _, comp := range entity.Components {
		if comp == component {
			continue
		}
		comp.OnNewComponent(entity, component)
	}
}

func (entity *Entity) RemoveComponent(component Component) {
	component.Destroy(entity)
	delete(entity.Components, component.Name())
}

func (entity *Entity) RemoveComponentName(name string) {
	if component, ok := entity.Components[name]; ok {
		component.Destroy(entity)
	}
	delete(entity.Components, name)
}

var components = make(map[string]reflect.Type)

func RegisterComponent(component Component) {
	components[component.Name()] = reflect.Indirect(reflect.ValueOf(component)).Type()
}

func NewComponent(name string) Component {
	compType, ok := components[name]
	if ok == false {
		log.Printf("Error loading component \"%v\", not registered!", name)
		return nil
	}

	component := reflect.New(compType)
	return component.Interface().(Component)
}
