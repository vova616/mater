package mater

import (
	"box2d"
)

type Entity struct {
	EntityClass
	Body *box2d.Body
	Enabled bool
}

type EntityClass interface {
	Update (dt float64)
}