package mater

import (
	. "box2d/vector2"
	"ftgl-go"
	"gl"
	"mater/render"
	"mater/util"
	"os"
)

var dbg = &util.Dbg

var TestFont *ftgl.Font
func init() {
	var err os.Error
	TestFont, err = ftgl.CreatePixmapFont("fonts/ttf-bitstream-vera-1.10/VeraMono.ttf")
	if err != nil {
		dbg.Printf("Error loading font")
		return
	}
	TestFont.SetFontFaceSize(72, 72);
}

type DebugData struct {
	SingleStep bool
	DebugView *DebugView
	Console *Console
}

func (dbg *DebugData) Init (mater *Mater) {
	dbg.Console = new(Console)
	dbg.Console.Init(mater)
}

func (mater *Mater) DebugDraw () {
	cam := mater.Scene.Camera
	gl.PushMatrix()
		gl.Color4f(0, 1, 0, .5)
		render.DrawCircle(Vector2{cam.ScreenSize.X / 2, cam.ScreenSize.Y / 2}, cam.ScreenSize.Y / 2.0 - 5.0, false)
		
		TestFont.RenderFont("TestText", ftgl.RENDER_ALL)

		//draw collision objects
		cam.PreDraw()
			mater.Dbg.DebugView.DrawDebugData()
		cam.PostDraw()

	gl.PopMatrix()
}