package main

import (
	. "github.com/teomat/mater"
	"github.com/teomat/mater/collision"
)

type CollisionCallback struct{}

func (cb *CollisionCallback) Name() string {
	return "CollisionCallback"
}

func (cb *CollisionCallback) Init(owner *Entity) {
	owner.Scene.Space.Callbacks.OnCollision = func(arb *collision.Arbiter) {
		
		//causes shapes to bounce off whenever they collide
		for i := 0; i < arb.NumContacts; i++ {
			c := &arb.Contacts[i]
			n := c.Normal
			n.Mult(-3)
			arb.ShapeA.Body.Velocity.Add(n)
		}
	}
}

func (cb *CollisionCallback) Update(owner *Entity, dt float64) {
	
}

func (cb *CollisionCallback) Destroy(owner *Entity) {

}

func (cb *CollisionCallback) Marshal(owner *Entity) ([]byte, error) {
	return ([]byte)("{}"), nil
}

func (cb *CollisionCallback) Unmarshal(owner *Entity, data []byte) error {
	return nil
}

func (cb *CollisionCallback) OnNewComponent(owner *Entity, other Component) {

}

func init() {
	RegisterComponent(&CollisionCallback{})
}
