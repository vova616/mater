package components

import (
	. "github.com/teomat/mater"
)

type CamTarget struct {}

func (ct *CamTarget) Name() string {
	return "CamTarget"
}

func (ct *CamTarget) Init(owner *Entity) {
	owner.Scene.Camera.Position = owner.Transform.Position
}

func (ct *CamTarget) Update(owner *Entity, dt float64) {
	owner.Scene.Camera.Position = owner.Transform.Position
}

func (ct *CamTarget) Destroy(owner *Entity) {
	
}

func (ct *CamTarget) Marshal(owner *Entity) ([]byte, error) {
	return ([]byte)("{}"), nil
}

func (ct *CamTarget) Unmarshal(owner *Entity, data []byte) error {
	return nil
}

func (ct *CamTarget) OnNewComponent(owner *Entity, other Component) {

}

func init() {
	RegisterComponent(&CamTarget{})
}
