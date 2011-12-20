package mater

import (
	"github.com/jteeuwen/glfw"
)

var QuickSavePath = "quicksave.json"

type OnKeyCallbackFunc func (mater *Mater, key, state int) OnKeyCallbackFunc

func (mater *Mater) OnKey (key, state int) {
	//first global key bindings
	if state == 1 {
		switch key {
			case glfw.KeyF1:
				println("Editmode!")
				return
		}
	}

	//
	mater.OnKeyCallback = mater.OnKeyCallback(mater,key, state)
	
}

func DefaultKeyCallback (mater *Mater, key, state int) OnKeyCallbackFunc {
	
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
			dbg.Printf("QuickSave\n")
			mater.SaveScene(QuickSavePath)
		case glfw.KeyF9:
			dbg.Printf("QuickLoad\n")
			mater.LoadScene(QuickSavePath)
			mater.Paused = true
		}
	}

	return DefaultKeyCallback
}