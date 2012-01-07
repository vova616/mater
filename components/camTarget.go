package components

import (
	. "mater"
)

type CamTarget struct {
	Target *Transform
}

func (ct *CamTarget) Name() string {
	return "CamTarget"
}

func (ct *CamTarget) Init(owner *Entity) {
	if xfComp, ok := owner.Components["Transform"]; ok {
		xf := xfComp.(*Transform)
		ct.Target = xf

		cam := owner.Scene.Camera
		cam.Position = ct.Target.Position
	}
}

func (ct *CamTarget) Update(owner *Entity, dt float64) {
	if ct.Target != nil && ct.Target.Transform != nil {
		owner.Scene.Camera.Position = ct.Target.Position
	}
}

func (ct *CamTarget) Destroy(owner *Entity) {
	ct.Target = nil
}

func (ct *CamTarget) Marshal(owner *Entity) ([]byte, error) {
	return ([]byte)("{}"), nil
}

func (ct *CamTarget) Unmarshal(owner *Entity, data []byte) error {
	return nil
}

func (ct *CamTarget) OnNewComponent(owner *Entity, other Component) {
	if other.Name() == "Transform" {
		xf := other.(*Transform)
		ct.Target = xf

		cam := owner.Scene.Camera
		cam.Position = ct.Target.Position
	}
}

func init() {
	RegisterComponent(&CamTarget{})
}
