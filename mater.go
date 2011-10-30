package mater

import (
	"gl"
	. "box2d/vector2"
	"gob"
	"os"
)

type Mater struct {
	DefaultCamera Camera
	Running, Paused bool
	Dbg DebugData
	Scene *Scene
	OnKeyCallback OnKeyCallbackFunc
}

func (mater *Mater) Init (cam *Camera) {
	mater.DefaultCamera = *cam
	dbg := &(mater.Dbg)
	dbg.Init(mater)
	mater.Scene = new(Scene)
	mater.Scene.Init(mater)
	mater.Scene.Camera = cam

	mater.Dbg.DebugView = NewDebugView(mater.Scene.World)

	mater.OnKeyCallback = DefaultKeyCallback
}

func (mater *Mater) OnResize (width, height int) {
	if height == 0 {
		height = 1
	}

	w, h := float64(width), float64(height)
	mater.DefaultCamera.ScreenSize = Vector2{w, h}
	mater.Scene.Camera.ScreenSize = Vector2{w, h}
	
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

var saveDirectory = "saves/"

func (mater *Mater) SaveScene (path string) os.Error{
	scene := mater.Scene

	path = saveDirectory + path

	file, err := os.Create(path)
	if err != nil {
		dbg.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(scene)
	if err != nil {
		dbg.Printf("Error encoding Scene: %v", err)
		return err
	}

	return nil
}

func (mater *Mater) LoadScene (path string) os.Error {

	var scene *Scene

	path = saveDirectory + path

	file, err := os.Open(path)
	if err != nil {
		dbg.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	scene = new(Scene)
	decoder := gob.NewDecoder(file)

	err = decoder.Decode(scene)
	if err != nil {
		dbg.Printf("Error decoding Scene: %v", err)
		return err
	}

	mater.Scene = scene
	scene.World.Enabled = true

	if mater.Scene.Camera == nil {
		cam := mater.DefaultCamera
		mater.Scene.Camera = &cam
	}

	mater.Dbg.DebugView.Reset(mater.Scene.World)

	return nil
}