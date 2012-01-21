package main

import (
	"github.com/teomat/mater/collision"
	"github.com/teomat/mater/engine"
)

type CollisionHandler struct{}

type CollisionCallback func(arb *collision.Arbiter)

func (cb *CollisionHandler) Name() string {
	return "CollisionHandler"
}

func (cb *CollisionHandler) Init(owner *engine.Entity) {
	owner.Scene.Space.Callbacks.OnCollision = func(arb *collision.Arbiter) {

		callCollisionFunc := func(s *collision.Shape) {
			if s.UserData != nil {
				//check the shape for a collision callback
				if ch, ok := s.UserData.(CollisionCallback); ok {
					ch(arb)
				} else if s.Body.UserData != nil {
					//check the shape's body for a collision callback
					if ch, ok := s.Body.UserData.(CollisionCallback); ok {
						ch(arb)
					}
				}
			} else if s.Body.UserData != nil {
				//check the shape's body for a collision callback
				if ch, ok := s.Body.UserData.(CollisionCallback); ok {
					ch(arb)
				}
			}
		}

		callCollisionFunc(arb.ShapeA)
		callCollisionFunc(arb.ShapeB)
	}
}

func (cb *CollisionHandler) Update(owner *engine.Entity, dt float64) {

}

func (cb *CollisionHandler) Destroy(owner *engine.Entity) {

}

func (cb *CollisionHandler) Marshal() ([]byte, error) {
	return ([]byte)("{}"), nil
}

func (cb *CollisionHandler) Unmarshal(data []byte) error {
	return nil
}

func (cb *CollisionHandler) OnNewComponent(owner *engine.Entity, other engine.Component) {

}

func init() {
	engine.RegisterComponent(&CollisionHandler{})
}
