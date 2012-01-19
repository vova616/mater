package main

import (
	. "github.com/teomat/mater"
	"github.com/teomat/mater/camera"
)

func OnNewComponent(entity *Entity, comp Component) {
	if comp.Name() == "Camera" {
		cam := comp.(*camera.Camera)

		if cam.IsMainCamera {
			MainCamera = cam
			MainCamera.ScreenSize = ScreenSize
		}
	}
}
