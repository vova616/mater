package transform

import (
	"github.com/teomat/mater/vect"
	"math"
)

type Rotation struct {
	//sine and cosine.
	S, C float64
}

func NewRotation(angle float64) Rotation {
	return Rotation{
		C: math.Cos(angle),
		S: math.Sin(angle),
	}
}

func (rot *Rotation) SetIdentity() {
	rot.S = 0
	rot.C = 1
}

func (rot *Rotation) SetAngle(angle float64) {
	rot.C = math.Cos(angle)
	rot.S = math.Sin(angle)
}

func (rot *Rotation) Angle() float64 {
	return math.Atan2(rot.S, rot.C)
}

//rotates the input vector.
func (rot *Rotation) RotateVect(v vect.Vect) vect.Vect {
	return vect.Vect{
		X: v.X*rot.C - v.Y*rot.S,
		Y: v.X*rot.S + v.Y*rot.C,
	}
}

func (rot *Rotation) RotateVectInv(v vect.Vect) vect.Vect {
	return vect.Vect{
		X: v.X*rot.C + v.Y*rot.S,
		Y: -v.X*rot.S + v.Y*rot.C,
	}
}

func RotateVect(v vect.Vect, r Rotation) vect.Vect {
	return r.RotateVect(v)
}

type Transform struct {
	Position vect.Vect
	Rotation
}

func NewTransform(pos vect.Vect, angle float64) Transform {
	return Transform{
		Position: pos,
		Rotation: NewRotation(angle),
	}
}

func (xf *Transform) SetIdentity() {
	xf.Position = vect.Vect{}
	xf.Rotation.SetIdentity()
}

func (xf *Transform) Set(pos vect.Vect, rot float64) {
	xf.Position = pos
	xf.SetAngle(rot)
}

//moves and roates the input vector.
func (xf *Transform) TransformVect(v vect.Vect) vect.Vect {
	return vect.Add(xf.Position, xf.RotateVect(v))
}

func (xf *Transform) TransformVectInv(v vect.Vect) vect.Vect {
	return vect.Add(vect.Mult(xf.Position, -1), xf.RotateVectInv(v))
}
