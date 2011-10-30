package texutil

import (
	"gl"
)

type Texture struct {
	gl.Texture
	Width, Height int
}