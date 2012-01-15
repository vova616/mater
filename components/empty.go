package components

import (
	. "github.com/teomat/mater"
)

//Can be embedded in other components to reduce boilerplate.
type Empty struct {}

func (empty Empty) Init(owner *Entity) {}

func (empty Empty) Update(owner *Entity, dt float64) {}

func (empty Empty) Destroy(owner *Entity) {}

func (empty Empty) Marshal(owner *Entity) ([]byte, error) {
	return ([]byte)("{}"), nil
}

func (empty Empty) Unmarshal(owner *Entity, data []byte) error {
	return nil
}

func (empty Empty) OnNewComponent(owner *Entity, other Component) {}
