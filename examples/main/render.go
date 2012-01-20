package main

import (
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/teomat/ftgl-go"
	"github.com/teomat/mater/camera"
	"github.com/teomat/mater/engine"
	"github.com/teomat/mater/render"
	"github.com/teomat/mater/vect"
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

func Draw(scene *engine.Scene) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.PushMatrix()
	defer gl.PopMatrix()

	gl.Color4f(0, 1, 0, .5)
	render.DrawCircle(vect.Vect{camera.ScreenSize.X / 2, camera.ScreenSize.Y / 2}, camera.ScreenSize.Y/2.0-5.0, false)

	if Settings.Paused {
		gl.Color3f(1, 1, 1)
		RenderFontAt("Paused", 20, 30)
	}

	cam := MainCamera
	if cam == nil {
		return
	}

	//draw collision objects
	cam.PreDraw()

	DrawDebugData(scene.Space)

	cam.PostDraw()
}
