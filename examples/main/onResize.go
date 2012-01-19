package main

import (
	"github.com/teomat/mater/vect"
	"github.com/teomat/mater/camera"
	"github.com/banthar/Go-OpenGL/gl"
)

var ScreenSize vect.Vect

func OnResize(width, height int) {
	if height == 0 {
		height = 1
	}

	w, h := float64(width), float64(height)
	ScreenSize = vect.Vect{w, h}
	camera.ScreenSize = ScreenSize
	if MainCamera != nil {
		MainCamera.ScreenSize = ScreenSize
	}

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	//camera centered at (0,0)
	gl.Ortho(0, w, h, 0, 1, -1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.Translated(.375, .375, 0)
}
