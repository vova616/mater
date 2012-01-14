package main

import (
	. "github.com/teomat/mater"
	"github.com/teomat/mater/collision"
	"github.com/teomat/mater/components"
	"github.com/teomat/mater/vect"
)

type CollisionCallbackTest struct{}

func (cht *CollisionCallbackTest) Name() string {
	return "CollisionCallbackTest"
}

func (cht *CollisionCallbackTest) Init(owner *Entity) {
	if comp, ok := owner.Components["Body"]; ok {
		cht.addHandlerToBody(comp)
	}
}

func (cht *CollisionCallbackTest) Update(owner *Entity, dt float64) {
	
}

func (cht *CollisionCallbackTest) Destroy(owner *Entity) {

}

func (cht *CollisionCallbackTest) Marshal(owner *Entity) ([]byte, error) {
	return ([]byte)("{}"), nil
}

func (cht *CollisionCallbackTest) Unmarshal(owner *Entity, data []byte) error {
	return nil
}

func (cht *CollisionCallbackTest) OnNewComponent(owner *Entity, other Component) {
	if other.Name() == "Body" {
		cht.addHandlerToBody(other)
	}
}

func (cht *CollisionCallbackTest) addHandlerToBody(comp Component) {
	if bodyComp, ok := comp.(*components.Body); ok {
		body := bodyComp.Body

		body.UserData = CollisionCallback(func(arb *collision.Arbiter) {
			arb.ShapeA.Body.Velocity.Add(vect.Mult(arb.Contacts[0].Normal, -3))
		})
	}
}

func init() {
	RegisterComponent(&CollisionCallbackTest{})
}
