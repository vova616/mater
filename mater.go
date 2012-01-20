package mater

type Mater struct {
	Scene           *Scene
}

func (mater *Mater) Init() {
	mater.Scene = new(Scene)
	mater.Scene.Init()
}

func (mater *Mater) Update(dt float64) {
	mater.Scene.Update(dt)
}
