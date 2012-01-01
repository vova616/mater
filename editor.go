package mater

type Editor interface {
	//called when changing scene
	Reset(mater *Mater)
	//called once per frame to allow rendering
	Draw(mater *Mater)
	//
	KeyCallback() OnKeyCallbackFunc
}
