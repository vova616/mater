package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/glfw"
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/teomat/mater/vect"
	. "github.com/teomat/mater"
	"log"
	//importing so the components can register themselves
	_ "github.com/teomat/mater/components"
)

var flags = struct {
	startPaused, help bool
	dbg bool
	file string
	buildExamples bool
}{}

func init () {
	flag.BoolVar(&flags.help, 
		"help", false, "Shows command line flags and default values")
	flag.BoolVar(&flags.startPaused, 
		"paused", false, "Starts the game in a paused state")
	flag.StringVar(&flags.file,
		"file", "", "WorldFile to load on start")
	flag.BoolVar(&flags.dbg, 
		"dbg", false, "Enable debugmode")
	flag.BoolVar(&flags.buildExamples,
		"examples", false, "Recreate the examples in \"saves/examples/\"")
}


func main () {

	log.SetFlags(log.Lshortfile)

	//parse flags
	flag.Parse()
	if flags.help {
		flag.PrintDefaults()
		return
	}

	if flags.buildExamples {
		allExamples()
		return
	}

	if err := glfw.Init(); err != nil {
		log.Printf("Error initializing glfw: %v\n", err)
		return
	}
	defer glfw.Terminate()

	wx, wy := 800, 600

	cam := new(Camera)
	cam.ScreenSize = vect.Vect{float64(wx), float64(wy)}
	cam.Position = vect.Vect{0, 0}
	cam.Scale = vect.Vect{32, 32}
	cam.Rotation = 0

	mater := new(Mater)
	mater.DefaultCamera = *cam
	mater.Init()

	mater.Paused = flags.startPaused

	if flags.file != "" {
		err := mater.LoadScene(flags.file)
		mater.Paused = true
		if err != nil {
			panic(err)
		}
	}
	
	glfw.OpenWindowHint(glfw.WindowNoResize, 1)
	//glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 1)

	if err := glfw.OpenWindow(wx, wy, 8, 8, 8, 8, 0, 8, glfw.Windowed); err != nil {
		log.Printf("Error opening Window: %v\n", err)
		return
	}
	defer glfw.CloseWindow()

	if gl.Init() != 0 {
		panic("gl error")
	}

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle("mater test")
	glfw.SetWindowSizeCallback(func(w, h int){mater.OnResize(w, h)})
	glfw.SetKeyCallback(func(k, s int){mater.OnKey(k, s)})
	
	//init opengl
	{
		gl.ClearColor(0, 0, 0, 0)
		gl.Enable(gl.BLEND)
		//gl.BlendFunc (gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
		gl.Enable (gl.TEXTURE_2D)
	}	

	printFPS := false
	var time, lastTime float64
	time = glfw.Time()
	lastTime = time
	var now, dt float64
	var frameCount, updateFrameCount int
	var acc, updateAcc float64
	//fix timestep to given fps
	const expectedFps = 30.0
	const expectedFrameTime = 1.0 / expectedFps
	var fps, updateFps int
	mater.Running = true
	for mater.Running && 1 == glfw.WindowParam(glfw.Opened) {
		time = glfw.Time()
		now = time
		dt = now - lastTime
		lastTime = now
		
		frameCount++
		acc += dt

		//fix update rate
		updateAcc += dt

		//execute console commands
		select {
			case command := <- mater.Dbg.Console.Command:
				mater.Dbg.Console.ExecuteCommand(command)
			default:
		}

		for updateAcc >= expectedFrameTime {
			updateFrameCount++

			if acc > 1 {
				updateFps = updateFrameCount
				updateFrameCount = 0
			}

			if !mater.Paused || mater.Dbg.SingleStep {
				mater.Update(expectedFrameTime)
				if mater.Dbg.SingleStep {mater.Dbg.SingleStep = false}
			}

			updateAcc -= expectedFrameTime
		}

		mater.Draw()

		mater.DebugDraw()

		glfw.SwapBuffers()

		if acc > 1 {
			fps = frameCount
			frameCount = 0
			if !mater.Paused && printFPS {
				fmt.Printf("---\n")
				fmt.Printf("FPS: %v\n", fps)
				fmt.Printf("Update FPS: %v\n", updateFps)
				fmt.Printf("Average frametime: %v\n", acc / float64(fps))
				fmt.Printf("---\n")
			}
			acc -= 1
		}
	}
}
