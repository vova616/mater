package texutil

import (
	"errors"
	"fmt"
	"github.com/banthar/Go-OpenGL/gl"
	"image"
	"image/color"
	"image/png"
	"os"
)

//
func LoadPng(path string) (*Texture, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg, err := png.DecodeConfig(file)

	if err != nil {
		return nil, err
	}

	file.Close()

	file, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(file)

	if err != nil {
		return nil, err
	}

	file.Close()

	var format, dataType gl.GLenum
	var data []byte
	switch cfg.ColorModel {
	case color.NRGBAModel, color.RGBAModel:
		format = gl.RGBA8
		dataType = gl.UNSIGNED_BYTE
	case color.NRGBA64Model, color.RGBA64Model:
		format = gl.RGBA16
		dataType = gl.UNSIGNED_SHORT
	default:
		panic(cfg.ColorModel)
		return nil, errors.New("Data Format not supported!")
	}

	switch cfg.ColorModel {
	case color.NRGBAModel:
		data = img.(*image.NRGBA).Pix
	case color.RGBAModel:
		data = img.(*image.RGBA).Pix
	case color.NRGBA64Model:
		data = img.(*image.NRGBA64).Pix
	case color.RGBA64Model:
		data = img.(*image.RGBA).Pix
	}

	texture := new(Texture)

	texture.Width = cfg.Width
	texture.Height = cfg.Height

	texture.Texture = gl.GenTexture()

	texture.Bind(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, cfg.Width, cfg.Height, 0, format, dataType, data)

	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		errString := gl.GetString(errCode)
		fmt.Printf("OpenGL Error: %v\n", errString)
	}

	texture.Unbind(gl.TEXTURE_2D)

	return texture, nil
}
