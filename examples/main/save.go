package main

import (
	"github.com/teomat/mater/engine"
	"bytes"
	"encoding/json"
	"log"
	"os"
)

var SaveDirectory = "saves/"

func saveScene(path string) error {
	path = Settings.SaveDir + path

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	dataString, err := json.MarshalIndent(scene, "", "\t")
	if err != nil {
		log.Printf("Error encoding Scene: %v", err)
		return err
	}

	buf := bytes.NewBuffer(dataString)
	n, err := buf.WriteTo(file)
	if err != nil {
		log.Printf("Error after writing %v characters to File: %v", n, err)
		return err
	}

	return nil
}

func loadScene(path string) error {

	var newScene *engine.Scene

	path = Settings.SaveDir + path

	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	newScene = new(engine.Scene)
	newScene.Callbacks = callbacks
	decoder := json.NewDecoder(file)

	err = decoder.Decode(newScene)
	if err != nil {
		log.Printf("Error decoding Scene: %v", err)
		return err
	}

	scene = newScene
	scene.Space.Enabled = true

	return nil
}
