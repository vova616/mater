package components

import (
	. "mater"
	"mater/collision"
	"json"
	"os"
)

type Body struct {
	*collision.Body
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

		body.Transform = *xf.Transform

		//set the transform to point to this one
		xf.Transform = &body.Transform
	}
	owner.Scene.Space.AddBody(body.Body)
}

func (body *Body) Update (owner *Entity, dt float64) {

}

func (body *Body) Destroy (owner *Entity) {
	if body.Body == nil {
		dbg.Printf("Error: Body component is not initialized correctly!")
		return
	}
	body.Body.Enabled = false
	owner.Scene.Space.RemoveBody(body.Body)
	body.Body.UserData = nil
}

func (body *Body) Marshal(owner *Entity) ([]byte, os.Error) {
	return json.Marshal(body.Body)
}

func (body *Body) Unmarshal(owner *Entity, data []byte) (os.Error) {
	if body.Body == nil {
		body.Body = collision.NewBody(collision.BodyType_Static)
	}
	err := json.Unmarshal(data, body.Body)

	if err != nil {
		return err
	}

	if body.Body == nil {
		return os.NewError("nil Body")
	}

	owner.Scene.Space.AddBody(body.Body)

	return nil
}

func init() {
	RegisterComponent(&Body{})
}
