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
	"mater/vect"
	"gl"
	. "mater/texutil"
	"math"
)

func DrawQuad (upperLeft, lowerRight vect.Vect, filled bool) {
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
	deg2grad = math.Pi / 180
)
func DrawCircle (pos vect.Vect, radius float64, filled bool) {
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
		gl.Vertex2d(pos.X + math.Cos(d) * radius, pos.Y + math.Sin(d) * radius)
	}
}

func DrawLine (start, end vect.Vect) {
	gl.Begin(gl.LINES)
	defer gl.End()
	
	gl.Vertex2d(start.X, start.Y)
	gl.Vertex2d(end.X, end.Y)
}

func DrawPoly (vertices []vect.Vect, vertCount int, filled bool) {
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

func DrawTextureQuad (texture *Texture, vertices [4]vect.Vect, texCoords [4]vect.Vect) {
	
	gl.BindTexture (gl.TEXTURE_2D, uint(texture.Texture))
	gl.Begin (gl.QUADS)

	for i := 0; i < 4; i++ {
		v := vertices[i]
		tc := texCoords[i]

		gl.TexCoord2d(tc.X, tc.Y)
		gl.Vertex2d(v.X, v.Y)
	}

	gl.End ()
}