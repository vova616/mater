package mater

import (
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/teomat/mater/vect"
	"math"
)

type Camera struct {
	ScreenSize, Position, Scale vect.Vect
	Rotation                    float64
}

func (cam Camera) PreDraw() {
	gl.PushMatrix()
	gl.Translated(cam.ScreenSize.X/2, cam.ScreenSize.Y/2, 0)

	gl.Rotated(360-cam.Rotation, 0, 0, 1)
	gl.Scaled(cam.Scale.X, cam.Scale.Y, 1)

	gl.Translated(-cam.Position.X, -cam.Position.Y, 0)
}

func (cam Camera) PostDraw() {
	gl.LoadIdentity()
	gl.PopMatrix()
}

func (cam Camera) Move(delta vect.Vect) {
	cam.Position = vect.Add(cam.Position, delta)
}

func (cam Camera) WorldToScreen(worldPos vect.Vect) vect.Vect {
	c := math.Cos(-cam.Rotation)
	s := math.Sin(-cam.Rotation)

	tx := worldPos.X - cam.Position.X
	tx += (cam.ScreenSize.X / 2) * cam.Scale.X

	ty := worldPos.Y - cam.Position.Y
	ty += (cam.ScreenSize.Y / 2) * cam.Scale.Y

	sx := c*tx - s*ty
	sy := s*tx + c*ty

	return vect.Vect{sx * cam.Scale.X, sy * cam.Scale.Y}
}

func (cam Camera) ScreenToWorld(screenPos vect.Vect) vect.Vect {
	c := math.Cos(cam.Rotation)
	s := math.Sin(cam.Rotation)

	tx := screenPos.X / cam.Scale.X
	ty := screenPos.Y / cam.Scale.Y

	sx := c*tx - s*ty
	sy := s*tx + c*ty

	sx += cam.Position.X
	sy += cam.Position.Y

	sx -= (cam.ScreenSize.X / 2) / cam.Scale.X
	sy -= (cam.ScreenSize.Y / 2) / cam.Scale.Y

	return vect.Vect{sx, sy}
}
