package camera

import (
	. "github.com/teomat/mater"
	"encoding/json"
)

func (cam *Camera) Name() string {
	return "Camera"
}

func (cam *Camera) Init(owner *Entity) {
	if cam.FollowTarget {
		cam.Transform.Position = owner.Transform.Position
	}
}

func (cam *Camera) Update(owner *Entity, dt float64) {
	if cam.FollowTarget {
		cam.Transform.Position = owner.Transform.Position
	}
}

func (cam *Camera) Marshal(owner *Entity) ([]byte, error) {
	return json.Marshal(cam)
}

func (cam *Camera) Unmarshal(owner *Entity, data []byte) error {
	cam.ScreenSize = ScreenSize
	err := json.Unmarshal(data, cam)

	if err != nil {
		return err
	}

	if cam.IsMainCamera {
		MainCamera = cam
	}
	return nil
}

func (cam *Camera) Destroy(owner *Entity) {}

func (cam *Camera) OnNewComponent(owner *Entity, other Component) {}

func init() {
	RegisterComponent(&Camera{})
}
