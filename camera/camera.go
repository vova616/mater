package camera

import (
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
)

var ScreenSize vect.Vect

type Camera struct {
	ScreenSize   vect.Vect `json:"-"`
	Scale        vect.Vect
	Transform    transform.Transform
	FollowTarget bool
	IsMainCamera bool
}

func (cam *Camera) PreDraw() {
	gl.PushMatrix()
	gl.Translated(cam.ScreenSize.X/2, cam.ScreenSize.Y/2, 0)

	gl.Rotated(360-cam.Transform.Angle(), 0, 0, 1)
	gl.Scaled(cam.Scale.X, cam.Scale.Y, 1)

	gl.Translated(-cam.Transform.Position.X, -cam.Transform.Position.Y, 0)
}

func (cam *Camera) PostDraw() {
	gl.LoadIdentity()
	gl.PopMatrix()
}

func (cam *Camera) Move(delta vect.Vect) {
	cam.Transform.Position = vect.Add(cam.Transform.Position, delta)
}

func (cam *Camera) WorldToScreen(worldPos vect.Vect) vect.Vect {
	c := cam.Transform.C
	s := cam.Transform.S

	tx := worldPos.X - cam.Transform.Position.X
	tx += (cam.ScreenSize.X / 2) * cam.Scale.X

	ty := worldPos.Y - cam.Transform.Position.Y
	ty += (cam.ScreenSize.Y / 2) * cam.Scale.Y

	sx := c*tx - s*ty
	sy := s*tx + c*ty

	return vect.Vect{sx * cam.Scale.X, sy * cam.Scale.Y}
}

func (cam *Camera) ScreenToWorld(screenPos vect.Vect) vect.Vect {
	c := cam.Transform.C
	s := cam.Transform.S

	tx := screenPos.X / cam.Scale.X
	ty := screenPos.Y / cam.Scale.Y

	sx := c*tx - s*ty
	sy := s*tx + c*ty

	sx += cam.Transform.Position.X
	sy += cam.Transform.Position.Y

	sx -= (cam.ScreenSize.X / 2) / cam.Scale.X
	sy -= (cam.ScreenSize.Y / 2) / cam.Scale.Y

	return vect.Vect{sx, sy}
}
