package main

import (
	"github.com/jteeuwen/glfw"
	. "github.com/teomat/mater"
	"log"
)

var QuickSavePath = "quicksave.json"

func OnKey(mater *Mater, key, state int) {
	if state == 1 {
		switch key {
		case 'P':
			Settings.Paused = !Settings.Paused
		case 'S':
			Settings.SingleStep = Settings.SingleStep
		//Escape
		case glfw.KeyEsc:
			Settings.Running = false
		case glfw.KeyF5:
			log.Printf("QuickSave\n")
			saveScene(mater.Scene, QuickSavePath)
		case glfw.KeyF9:
			log.Printf("QuickLoad\n")
			loadScene(mater.Scene, QuickSavePath)
			Settings.Paused = true
		}
	}
}
