/*
* Copyright (c) 2011 Matteo Goggi
*
* This software is provided 'as-is', without any express or implied 
* warranty.  In no event will the authors be held liable for any damages 
* arising from the use of this software. 
* Permission is granted to anyone to use this software for any purpose, 
* including commercial applications, and to alter it and redistribute it 
* freely, subject to the following restrictions: 
* 1. The origin of this software must not be misrepresented; you must not 
* claim that you wrote the original software. If you use this software 
* in a product, an acknowledgment in the product documentation would be 
* appreciated but is not required. 
* 2. Altered source versions must be plainly marked as such, and must not be 
* misrepresented as being the original software. 
* 3. This notice may not be removed or altered from any source distribution. 
*/
package render

import (
	"ftgl-go"
	"gl"
	"mater/log"
	"os"
)

var dbg = &log.Dbg

var Font *ftgl.Font
func init() {
	var err os.Error
	Font, err = ftgl.CreatePixmapFont("fonts/ttf-bitstream-vera-1.10/VeraMono.ttf")
	if err != nil {
		dbg.Printf("Error loading main font:, %v", err)
		return
	}
	Font.SetFontFaceSize(20, 20);
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