package mater

import (
	. "box2d/vector2"
	"gl"
)

type Mater struct {
	DefaultCamera Camera
	ScreenSize Vector2
	Running, Paused bool
	Dbg DebugData
	Scene *Scene
	OnKeyCallback OnKeyCallbackFunc
}

func (mater *Mater) Init () {
	dbg := &(mater.Dbg)
	dbg.Init(mater)
	mater.Scene = new(Scene)
	mater.Scene.Init(mater)

	cl := &MaterContactListener{mater}
	mater.Scene.World.SetContactListener(cl)
	mater.Scene.World.SetContactFilter(cl)

	if dbg.DebugView == nil {
		mater.Dbg.DebugView = NewDebugView(mater.Scene.World)
	} else {
		mater.Dbg.DebugView.Reset(mater.Scene.World)
	}

	mater.OnKeyCallback = DefaultKeyCallback
}

func (mater *Mater) OnResize (width, height int) {
	if height == 0 {
		height = 1
	}

	w, h := float64(width), float64(height)
	mater.ScreenSize = Vector2{w, h}
	mater.DefaultCamera.ScreenSize = mater.ScreenSize
	if mater.Scene != nil {
		mater.Scene.Camera.ScreenSize = mater.ScreenSize
	}
	
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	//camera centered at (0,0)
	gl.Ortho(0, w, h, 0, 1, -1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.Translated(.375, .375, 0)
}

func (mater *Mater) Update (dt float64) {
	
	mater.Scene.Update(dt)
}

func (mater *Mater) Draw () {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	mater.Scene.Camera.PreDraw()
	{

	}
	mater.Scene.Camera.PostDraw()
}
