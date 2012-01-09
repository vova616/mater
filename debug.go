package mater

import (
	"mater/vect"
	"gl"
	"mater/render"
)

type DebugData struct {
	SingleStep bool
	DebugView *DebugView
	Console *Console
}

func (dd *DebugData) Init (mater *Mater) {
	dd.Console = new(Console)
	dd.Console.Init(mater)
}

func (mater *Mater) DebugDraw () {
	cam := mater.Scene.Camera
	gl.PushMatrix()
		gl.Color4f(0, 1, 0, .5)
		render.DrawCircle(vect.Vect{cam.ScreenSize.X / 2, cam.ScreenSize.Y / 2}, cam.ScreenSize.Y / 2.0 - 5.0, false)
		
		if mater.Paused {
			gl.Color3f(1, 1, 1)
			render.RenderFontAt("Paused", 20, 30)
		}
		

		//draw collision objects
		cam.PreDraw()
			mater.Dbg.DebugView.DrawDebugData()
		cam.PostDraw()

	gl.PopMatrix()
}