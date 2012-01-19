package mater

type Mater struct {
	Running, Paused bool
	Dbg             DebugData
	Scene           *Scene
	OnKeyCallback   OnKeyCallbackFunc
}

func (mater *Mater) Init() {
	dbg := &(mater.Dbg)
	dbg.Init(mater)
	mater.Scene = new(Scene)
	mater.Scene.Init(mater)

	mater.OnKeyCallback = DefaultKeyCallback
}

func (mater *Mater) Update(dt float64) {
	mater.Scene.Update(dt)
}
