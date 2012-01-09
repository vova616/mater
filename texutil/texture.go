package texutil

import (
	"github.com/banthar/Go-OpenGL/gl"
)

type Texture struct {
	gl.Texture
	Width, Height int
}