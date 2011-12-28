package components

import (
	. "mater"
	"box2d"
	"json"
	"os"
)

type Body struct {
	*box2d.Body
}

func (body *Body) Name () string {
	return "Body"
}

func (body *Body) Init (owner *Entity) {
	if body.Body == nil {
		dbg.Printf("Error: Body component is not initialized correctly!")
		return
	}

	body.Body.UserData = owner

	if tcomp, ok := owner.Components["Transform"]; ok {
		//the owner already has a transform attached, change this bodies transform to it
		xf := tcomp.(*Transform)

		body.SetTransform(&xf.Position, xf.Angle())

		//set the transform to point to this one
		xf.Transform = body.Transform()
	}
	body.RegisterBody(owner.Scene.World)
}

func (body *Body) Update (owner *Entity, dt float64) {

}

func (body *Body) Destroy (owner *Entity) {
	if body.Body == nil {
		dbg.Printf("Error: Body component is not initialized correctly!")
		return
	}
	owner.Scene.World.RemoveBody(body.Body)
	body.Body.SetEnabled(false)
	body.Body.UserData = nil
}

func (body *Body) Marshal(owner *Entity) ([]byte, os.Error) {
	return json.Marshal(body.Body)
}

func (body *Body) Unmarshal(owner *Entity, data []byte) (os.Error) {
	if body.Body == nil {
		body.Body = new(box2d.Body)
	}
	err := json.Unmarshal(data, body.Body)

	if err != nil {
		return err
	}

	if body.Body == nil {
		return os.NewError("nil Body")
	}

	return nil
}

func init() {
	RegisterComponent(&Body{})
}
