package mater

import (
	. "box2d"
)

type MaterContactListener struct {
	Mater *Mater
}

func (m *MaterContactListener) BeginContact (contact *Contact) {

}

func (m *MaterContactListener) EndContact (contact *Contact) {

}

func (m *MaterContactListener) PreSolve (contact *Contact, oldManifold *Manifold) {

}

func (m *MaterContactListener) PostSolve (contact *Contact, impulse *ContactConstraint) {

}

func (m *MaterContactListener) ShouldCollide (fixtureA, fixtureB *Fixture) bool {
	
	return true
}
