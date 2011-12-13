package mater

import (
	"bytes"
	"json"
	"os"
)

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

	//encoder := json.NewEncoder(file)
	//err = encoder.Encode(scene)

	dataString, err := json.MarshalIndent(scene, "", "\t")
	if err != nil {
		dbg.Printf("Error encoding Scene: %v", err)
		return err
	}

	buf := bytes.NewBuffer(dataString)
	n, err := buf.WriteTo(file)
	if err != nil {
		dbg.Printf("Error after writing %v characters to File: %v", n, err)
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
	decoder := json.NewDecoder(file)

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
	} else {
		mater.Scene.Camera.ScreenSize = mater.ScreenSize
	}

	mater.Dbg.DebugView.Reset(mater.Scene.World)

	return nil
}