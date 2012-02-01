package main

import (
	"github.com/jteeuwen/glfw"
	"log"
)

var QuickSavePath = "quicksave.json"

func OnKey(key, state int) {
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
			saveSpace(QuickSavePath)
		case glfw.KeyF9:
			log.Printf("QuickLoad\n")
			loadSpace(QuickSavePath)
			Settings.Paused = true
		}
	}
}
