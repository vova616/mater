package mater

import (
	"github.com/jteeuwen/glfw"
)

var QuickSavePath = "saves/quicksave.dat"

func (mater *Mater) OnKey (key, state int) {
	if state == 1 {
		switch key {
		case glfw.KeyF3:
			mater.Dbg.DrawDebugGraph = !mater.Dbg.DrawDebugGraph
		case 'P':
			mater.Paused = !mater.Paused
		case 'S':
			mater.Dbg.SingleStep = !mater.Dbg.SingleStep
		//Escape
		case glfw.KeyEsc:
			mater.Running = false
		case glfw.KeyF5:
			dbg.Printf("QuickSave\n")
			mater.SaveScene(QuickSavePath)
		case glfw.KeyF9:
			dbg.Printf("QuickLoad\n")
			mater.LoadScene(QuickSavePath)
			mater.Paused = true
		}
	}
}