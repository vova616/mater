package camera

import (
	"github.com/teomat/mater/engine"
	"encoding/json"
)

func (cam *Camera) Name() string {
	return "Camera"
}

func (cam *Camera) Init(owner *engine.Entity) {
	if cam.FollowTarget {
		cam.Transform.Position = owner.Transform.Position
	}
}

func (cam *Camera) Update(owner *engine.Entity, dt float64) {
	if cam.FollowTarget {
		cam.Transform.Position = owner.Transform.Position
	}
}

func (cam *Camera) Marshal(owner *engine.Entity) ([]byte, error) {
	return json.Marshal(cam)
}

func (cam *Camera) Unmarshal(owner *engine.Entity, data []byte) error {
	cam.ScreenSize = ScreenSize
	err := json.Unmarshal(data, cam)

	if err != nil {
		return err
	}

	return nil
}

func (cam *Camera) Destroy(owner *engine.Entity) {}

func (cam *Camera) OnNewComponent(owner *engine.Entity, other engine.Component) {}

func init() {
	engine.RegisterComponent(&Camera{})
}
