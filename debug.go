package mater

import (
	"os"
	"gl"
	. "box2d/vector2"
	"mater/util"
	. "mater/texutil"
	"ftgl-go"
)

var dbg = &util.Dbg
var TestTexture *Texture

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
	DrawDebugGraph bool
	TimeData struct {
		UpdateTime []float64
		UpdateTimeIndex int
		DrawTime []float64
		DrawTimeIndex int
		Values int
	}
	DebugView *DebugView
	Console *Console
}

func (dbg *DebugData) Init (mater *Mater) {
	dbg.TimeData.Values = 128
	dbg.TimeData.DrawTime = make([]float64, dbg.TimeData.Values)
	dbg.TimeData.UpdateTime = make([]float64, dbg.TimeData.Values)
	dbg.Console = new(Console)
	dbg.Console.Init(mater)
}

func (mater *Mater) DebugDraw () {
	dbg := &(mater.Dbg)
	cam := mater.Scene.Camera
	gl.PushMatrix()
		gl.Color4f(0, 1, 0, .5)
		Render.DrawCircle(Vector2{cam.ScreenSize.X / 2, cam.ScreenSize.Y / 2}, cam.ScreenSize.Y / 2.0 - 5.0, false)
		
		TestFont.RenderFont("TestText", ftgl.RENDER_ALL)
		/*{
			size := 64.0
			verts := [4]Vector2{
				Vector2{-size, -size},
				Vector2{size, -size},
				Vector2{size, size},
				Vector2{-size, size},
			}

			tsize := 1.0
			tc := [4]Vector2{
				Vector2{-tsize, -tsize},
				Vector2{tsize, -tsize},
				Vector2{tsize, tsize},
				Vector2{-tsize, tsize},
			}
			Render.DrawTextureQuad(TestTexture, verts, tc)
		}*/

		//draw collision objects
		cam.PreDraw()
			mater.Dbg.DebugView.DrawDebugData()
		cam.PostDraw()

		if mater.Dbg.DrawDebugGraph {
			//Draw Time graphs
			{
				//fill background
				gl.Color3f(1, 1, 1)
				const xScale = 2.0
				sx := xScale * float64(dbg.TimeData.Values)
				sy := 128.0

				gl.Begin(gl.QUADS)
					gl.Vertex2d(0, 0)
					gl.Vertex2d(0, sy)
					gl.Vertex2d(sx - xScale, sy)
					gl.Vertex2d(sx - xScale, 0)
				gl.End()

				//1 line : 0.02 ms
				const yRuler = 10.0
				const yScale = 5000.0
				gl.Color3f(0, 1, 0)
				gl.Begin(gl.LINES)
					for i := 0.0; i <= sy; i += yRuler {
						gl.Vertex2d(0.0, sy - i)
						gl.Vertex2d(sx - xScale, sy - i)
					}
					const k_160 = 1.0 / 60.0
					gl.Color3f(1, 0, 1)
					gl.Vertex2d(0.0, sy - k_160 * yScale)
					gl.Vertex2d(sx - xScale, sy - k_160 * yScale)
				gl.End()

				//update time
				gl.Color3f(1, 0, 0)
				gl.Begin(gl.LINE_LOOP)
				gl.Vertex2d(0, sy)
					for i, v := range(mater.Dbg.TimeData.UpdateTime) {
						gl.Vertex2d(float64(i) * xScale, sy - v * yScale)
					}
				gl.Vertex2d(sx - xScale, sy)
				gl.End()

				//draw time
				gl.Color3f(0, 0, 1)
				gl.Begin(gl.LINE_LOOP)
				gl.Vertex2d(0, sy)
					for i, v := range(mater.Dbg.TimeData.DrawTime) {
						gl.Vertex2d(float64(i) * xScale, sy - v * yScale)
					}
				gl.Vertex2d(sx - xScale, sy)
				gl.End()
			}
		}
	gl.PopMatrix()
}