package mater

import (
	"gl"
	. "box2d/vector2"
	"gob"
	"os"
	"json"
	"io/ioutil"
)

type Mater struct {
	Camera *Camera
	Running, Paused bool
	Dbg DebugData
	Scene *Scene
}

func (mater *Mater) Init (cam *Camera) {
	mater.Camera = cam
	dbg := &(mater.Dbg)
	dbg.Init()
	mater.Scene = new(Scene)
	mater.Scene.Init()

	mater.Dbg.DebugView = NewDebugView(mater.Scene.World)
}

func (mater *Mater) OnResize (width, height int) {
	if height == 0 {
		height = 1
	}

	w, h := float64(width), float64(height)
	mater.Camera.ScreenSize = Vector2{w, h}
	
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
	mater.Camera.PreDraw()
	{

	}
	mater.Camera.PostDraw()
}

var encodeType = "gob"
func (mater *Mater) SaveScene (path string) os.Error{
	scene := mater.Scene
	if encodeType == "gob" {

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
	} else if encodeType == "json" {

		sceneData, err := json.MarshalIndent(scene, "", "\t")
		if err != nil {
			dbg.Printf("Error marshaling World: %v", err)
			return err
		}

		file, err := os.Create(path)
		if err != nil {
			dbg.Printf("Error opening File: %v", err)
			return err
		}
		defer file.Close()

		if _, err := file.Write(sceneData); err != nil {
			dbg.Printf("Error writing File: %v", err)
			return err
		}

	}

	return os.NewError("Unknown encoding")
}

func (mater *Mater) LoadScene (path string) os.Error {

	var scene *Scene

	if encodeType == "gob" {

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

	} else if encodeType == "json" {
		file, err := os.Open(path)
		if err != nil {
			dbg.Printf("Error opening File: %v", err)
			return err
		}
		defer file.Close()

		var data []byte
		data, err = ioutil.ReadAll(file)
		if err != nil {
			dbg.Printf("Error reading File: %v", err)
			return err
		}

		scene = new(Scene)
		err = json.Unmarshal(data, scene)

		if err != nil {
			dbg.Printf("Error unmarshaling World: %v", err)
			return err
		}
	}

	mater.Scene = scene
	scene.World.Enabled = true

	mater.Dbg.DebugView.Reset(mater.Scene.World)

	return nil
}