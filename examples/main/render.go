package main

import (
	. "github.com/teomat/mater"
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/teomat/mater/render"
	"github.com/teomat/mater/vect"
	"github.com/teomat/ftgl-go"
	"log"
)

var Font *ftgl.Font

func init() {
	var err error
	Font, err = ftgl.CreatePixmapFont("fonts/ttf-bitstream-vera-1.10/VeraMono.ttf")
	if err != nil {
		log.Printf("Error loading main font:, %v", err)
		panic(err)
	}
	Font.SetFontFaceSize(20, 20)
}

func RenderFontAt(text string, x, y float64) {
	gl.PushMatrix()

	gl.RasterPos2d(x, y)

	Font.RenderFont(text, ftgl.RENDER_ALL)

	gl.PopMatrix()
}

func Draw(mater *Mater) {

	gl.Clear(gl.COLOR_BUFFER_BIT)

	cam := mater.Scene.Camera
	gl.PushMatrix()
	gl.Color4f(0, 1, 0, .5)
	render.DrawCircle(vect.Vect{cam.ScreenSize.X / 2, cam.ScreenSize.Y / 2}, cam.ScreenSize.Y/2.0-5.0, false)

	if mater.Paused {
		gl.Color3f(1, 1, 1)
		RenderFontAt("Paused", 20, 30)
	}

	//draw collision objects
	cam.PreDraw()
	
	DrawDebugData(mater.Scene.Space)
	
	cam.PostDraw()

	gl.PopMatrix()
}
