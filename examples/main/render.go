package main

import (
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/teomat/ftgl-go"
	"github.com/teomat/mater/vect"
	"log"
	"math"
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

func DrawQuad(upperLeft, lowerRight vect.Vect, filled bool) {
	if filled {
		gl.Begin(gl.QUADS)
	} else {
		gl.Begin(gl.LINE_LOOP)
	}
	defer gl.End()

	gl.Vertex2d(upperLeft.X, upperLeft.Y)
	gl.Vertex2d(upperLeft.X, lowerRight.Y)
	gl.Vertex2d(lowerRight.X, lowerRight.Y)
	gl.Vertex2d(lowerRight.X, upperLeft.Y)
}

const (
	circlestep = 45
	deg2grad   = math.Pi / 180
)

func DrawCircle(pos vect.Vect, radius float64, filled bool) {
	if filled {
		gl.Begin(gl.TRIANGLE_FAN)
		gl.Vertex2d(pos.X, pos.Y)
	} else {
		gl.Begin(gl.LINE_LOOP)
	}
	defer gl.End()

	var d float64
	for i := 0.0; i <= 360; i += circlestep {
		d = deg2grad * i
		gl.Vertex2d(pos.X+math.Cos(d)*radius, pos.Y+math.Sin(d)*radius)
	}
}

func DrawLine(start, end vect.Vect) {
	gl.Begin(gl.LINES)
	defer gl.End()

	gl.Vertex2d(start.X, start.Y)
	gl.Vertex2d(end.X, end.Y)
}

func DrawPoly(vertices []vect.Vect, vertCount int, filled bool) {
	if filled {
		gl.Begin(gl.TRIANGLE_FAN)
		gl.Vertex2d(vertices[0].X, vertices[0].Y)
	} else {
		gl.Begin(gl.LINE_LOOP)
	}
	defer gl.End()

	for i := 0; i < vertCount; i++ {
		v := vertices[i]
		gl.Vertex2d(v.X, v.Y)
	}
}

func Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.PushMatrix()
	defer gl.PopMatrix()

	gl.Color4f(0, 1, 0, .5)
	DrawCircle(vect.Vect{ScreenSize.X / 2, ScreenSize.Y / 2}, ScreenSize.Y/2.0-5.0, false)

	if Settings.Paused {
		gl.Color3f(1, 1, 1)
		RenderFontAt("Paused", 20, 30)
	}

	//draw collision objects

	gl.PushMatrix()
	
	gl.Translated(ScreenSize.X / 2, ScreenSize.Y / 2, 0)
	gl.Scaled(32, 32, 0)

	DrawDebugData(space)

	gl.PopMatrix()
}
