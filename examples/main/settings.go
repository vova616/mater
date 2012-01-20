package main

import (
	"os"
	"encoding/json"
	. "github.com/teomat/mater"
	"log"
	"bytes"
	"github.com/jteeuwen/glfw"
)

const settingsPath = "settings.json"

var Settings = struct {
	Resolution struct {
		Width int
		Height int
	}
	SaveDir string
	Running bool `json:"-"`
	Paused bool `json:"-"`
	SingleStep bool `json:"-"`
}{
	Resolution: struct{
		Width int
		Height int
	}{
		Width: 800,
		Height: 600,
	},
	SaveDir: "saves/",
}

func saveSettingsFile()  error {
	file, err := os.Create(settingsPath)
	if err != nil {
		log.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	dataString, err := json.MarshalIndent(&Settings, "", "\t")
	if err != nil {
		log.Printf("Error encoding Settings: %v", err)
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

func loadSettingsFile() error {
	file, err := os.Open(settingsPath)
	if err != nil {
		log.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&Settings)
	if err != nil {
		log.Printf("Error decoding Settings: %v", err)
		return err
	}

	return nil
}

func reloadSettings(scene *Scene) {
	glfw.SetWindowSize(Settings.Resolution.Width, Settings.Resolution.Height)
}
