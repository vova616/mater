package main

import (
	"github.com/teomat/mater/engine"
	"github.com/teomat/mater/camera"
)

func OnNewComponent(entity *engine.Entity, comp engine.Component) {
	if comp.Name() == "Camera" {
		cam := comp.(*camera.Camera)

		if cam.IsMainCamera {
			MainCamera = cam
			MainCamera.ScreenSize = ScreenSize
		}
	}
}
