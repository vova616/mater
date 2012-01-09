package render

import (
	"ftgl-go"
	"gl"
	"log"
)

var Font *ftgl.Font

func init() {
	var err error
	Font, err = ftgl.CreatePixmapFont("fonts/ttf-bitstream-vera-1.10/VeraMono.ttf")
	if err != nil {
		log.Printf("Error loading main font:, %v", err)
		return
	}
	Font.SetFontFaceSize(20, 20)
}

func SetFontFaceSize(size, res uint) {
	Font.SetFontFaceSize(size, res)
}

func RenderFont(text string) {
	Font.RenderFont(text, ftgl.RENDER_ALL)
}

func RenderFontAt(text string, x, y float64) {
	gl.PushMatrix()

	gl.RasterPos2d(x, y)

	Font.RenderFont(text, ftgl.RENDER_ALL)

	gl.PopMatrix()
}
