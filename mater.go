package mater

type Mater struct {
	Running, Paused bool
	Dbg             DebugData
	Scene           *Scene
	Callbacks struct {
		OnNewComponent func(entity *Entity, comp Component)
	}
}

func (mater *Mater) Init() {
	mater.Scene = new(Scene)
	mater.Scene.Init(mater)
}

func (mater *Mater) Update(dt float64) {
	mater.Scene.Update(dt)
}
