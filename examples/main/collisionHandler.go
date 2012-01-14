package main

import (
	. "github.com/teomat/mater"
	"github.com/teomat/mater/collision"
)

type CollisionHandler struct{}

type CollisionCallback func(arb *collision.Arbiter)

func (cb *CollisionHandler) Name() string {
	return "CollisionHandler"
}

func (cb *CollisionHandler) Init(owner *Entity) {
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

func (cb *CollisionHandler) Update(owner *Entity, dt float64) {
	
}

func (cb *CollisionHandler) Destroy(owner *Entity) {

}

func (cb *CollisionHandler) Marshal(owner *Entity) ([]byte, error) {
	return ([]byte)("{}"), nil
}

func (cb *CollisionHandler) Unmarshal(owner *Entity, data []byte) error {
	return nil
}

func (cb *CollisionHandler) OnNewComponent(owner *Entity, other Component) {

}

func init() {
	RegisterComponent(&CollisionHandler{})
}
