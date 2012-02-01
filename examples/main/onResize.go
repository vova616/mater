package main

import (
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/teomat/mater/vect"
)

var ScreenSize vect.Vect

//Callback for window resize events.
//Updates Settings.ScreenSize as well as the size of the main camera.
func OnResize(width, height int) {
	if height == 0 {
		height = 1
	}

	w, h := float64(width), float64(height)
	ScreenSize = vect.Vect{w, h}

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	//camera centered at (0,0)
	gl.Ortho(0, w, h, 0, 1, -1)
	gl.MatrixMode(gl.MODELVIEW)
	//gl.Translated(.375, .375, 0)

	//gl.Translated(-w/8, -h/2, 0)
	//gl.Scaled(32, 32, 1)
}
