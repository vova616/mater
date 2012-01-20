package main

import (
	"github.com/teomat/mater/engine"
	"github.com/teomat/mater/collision"
	"github.com/teomat/mater/components"
	"github.com/teomat/mater/vect"
)

type CollisionCallbackTest struct{}

func (cht *CollisionCallbackTest) Name() string {
	return "CollisionCallbackTest"
}

func (cht *CollisionCallbackTest) Init(owner *engine.Entity) {
	if comp, ok := owner.Components["Body"]; ok {
		cht.addHandlerToBody(comp)
	}
}

func (cht *CollisionCallbackTest) Update(owner *engine.Entity, dt float64) {

}

func (cht *CollisionCallbackTest) Destroy(owner *engine.Entity) {

}

func (cht *CollisionCallbackTest) Marshal(owner *engine.Entity) ([]byte, error) {
	return ([]byte)("{}"), nil
}

func (cht *CollisionCallbackTest) Unmarshal(owner *engine.Entity, data []byte) error {
	return nil
}

func (cht *CollisionCallbackTest) OnNewComponent(owner *engine.Entity, other engine.Component) {
	if other.Name() == "Body" {
		cht.addHandlerToBody(other)
	}
}

func (cht *CollisionCallbackTest) addHandlerToBody(comp engine.Component) {
	if bodyComp, ok := comp.(*components.Body); ok {
		body := bodyComp.Body

		body.UserData = CollisionCallback(func(arb *collision.Arbiter) {
			arb.ShapeA.Body.Velocity.Add(vect.Mult(arb.Contacts[0].Normal, -3))
		})
	}
}

func init() {
	engine.RegisterComponent(&CollisionCallbackTest{})
}
