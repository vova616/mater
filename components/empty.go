package components

import (
	"github.com/teomat/mater/engine"
)

//Can be embedded in other components to reduce boilerplate.
type Empty struct{}

func (empty Empty) Init(owner *engine.Entity) {}

func (empty Empty) Update(owner *engine.Entity, dt float64) {}

func (empty Empty) Destroy(owner *engine.Entity) {}

func (empty Empty) Marshal() ([]byte, error) {
	return ([]byte)("{}"), nil
}

func (empty Empty) Unmarshal(data []byte) error {
	return nil
}

func (empty Empty) OnNewComponent(owner *engine.Entity, other engine.Component) {}
