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
			mater.Paused = !mater.Paused
		case 'S':
			mater.Dbg.SingleStep = !mater.Dbg.SingleStep
		//Escape
		case glfw.KeyEsc:
			mater.Running = false
		case glfw.KeyF5:
			log.Printf("QuickSave\n")
			mater.SaveScene(QuickSavePath)
		case glfw.KeyF9:
			log.Printf("QuickLoad\n")
			mater.LoadScene(QuickSavePath)
			mater.Paused = true
		}
	}
}
