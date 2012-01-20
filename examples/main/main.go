package main

import (
	"flag"
	"fmt"
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/jteeuwen/glfw"
	"github.com/teomat/mater/engine"
	"github.com/teomat/mater/vect"
	"log"
	//importing so the components can register themselves
	"github.com/teomat/mater/camera"
	_ "github.com/teomat/mater/components"
	"os"
	"runtime/pprof"
)

var flags = struct {
	startPaused   bool
	file          string
	buildExamples bool
	cpuprofile    string
}{}

func init() {
	flag.BoolVar(&flags.startPaused,
		"paused", false, "Start the game in a paused state.")
	flag.StringVar(&flags.file,
		"file", "", "Load the given savefile located in \"./saves/\" on start. (e.g. -file=\"quicksave.json\")")
	flag.BoolVar(&flags.buildExamples,
		"examples", false, "Rebuild example savefiles in \"./saves/examples/\".")
	flag.StringVar(&flags.cpuprofile,
		"cpuprofile", "", "Write cpu profile to file.")
}

var MainCamera *camera.Camera
var console Console
var callbacks engine.Callbacks

var scene *engine.Scene

func main() {
	log.SetFlags(log.Lshortfile)

	//parse flags
	flag.Parse()

	if flags.cpuprofile != "" {
		f, err := os.Create(flags.cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	loadSettingsFile()
	Settings.Paused = flags.startPaused

	if flags.buildExamples {
		allExamples()
		return
	}

	wx, wy := 800, 600

	//initialize opengl & glfw
	{
		//init glfw
		if err := glfw.Init(); err != nil {
			log.Printf("Error initializing glfw: %v\n", err)
			return
		}
		defer glfw.Terminate()

		//set window hints
		glfw.OpenWindowHint(glfw.WindowNoResize, 1)

		//create the window
		if err := glfw.OpenWindow(wx, wy, 8, 8, 8, 8, 0, 8, glfw.Windowed); err != nil {
			log.Printf("Error opening Window: %v\n", err)
			return
		}
		defer glfw.CloseWindow()

		//init opengl
		if gl.Init() != 0 {
			panic("gl error")
		}

		//glfw config
		{
			glfw.SetSwapInterval(1)
			glfw.SetWindowTitle("mater test")
		}

		//set additional opengl stuff
		{
			gl.ClearColor(0, 0, 0, 0)
			gl.Enable(gl.BLEND)
			//gl.BlendFunc (gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
			gl.Enable(gl.TEXTURE_2D)
		}
	}

	//setup scene related stuff
	{
		//setup default camera
		MainCamera = new(camera.Camera)
		cam := MainCamera
		camera.ScreenSize = vect.Vect{float64(wx), float64(wy)}
		cam.Transform.Position = vect.Vect{0, 0}
		cam.Scale = vect.Vect{32, 32}
		cam.Transform.SetAngle(0)

		//create empty scene
		scene = new(engine.Scene)
		scene.Init()
	}

	//reload settings so they take effect
	reloadSettings()

	//set callbacks
	{
		callbacks.OnNewComponent = onNewComponent
		scene.Callbacks = callbacks
		glfw.SetWindowSizeCallback(func(w, h int) { OnResize(w, h) })
		glfw.SetKeyCallback(func(k, s int) { OnKey(k, s) })
	}

	//init debug console
	console.Init()

	//load savefile passed from the commandline if any
	if flags.file != "" {
		err := loadScene(flags.file)
		Settings.Paused = true
		if err != nil {
			panic(err)
		}
	}

	//if set to true once a second 
	printFPS := false

	//fix timestep to given fps
	const expectedFps = 30.0
	const expectedFrameTime = 1.0 / expectedFps

	//the time at the start of the last frame
	lastTime := 0.0

	acc := 0.0
	updateAcc := 0.0

	frameCount := 0
	updateFrameCount := 0

	fps := 0
	updateFps := 0

	Settings.Running = true

	for Settings.Running && glfw.WindowParam(glfw.Opened) == 1 {
		time := glfw.Time()
		//get the time elapsed since the last frame
		dt := time - lastTime
		lastTime = time

		//advance framecount and accumulators
		frameCount++
		acc += dt
		updateAcc += dt

		//execute console commands if any
		select {
		case command := <-console.Command:
			console.ExecuteCommand(command)
		default:
		}

		//update the scene at a fixed timestep
		for updateAcc >= expectedFrameTime {
			updateFrameCount++

			//if one second has passed update the fps and reset the framecount
			if acc > 1 {
				updateFps = updateFrameCount
				updateFrameCount = 0
			}

			//only update if not paused or if set to advance a single frame
			if !Settings.Paused || Settings.SingleStep {
				scene.Update(expectedFrameTime)
				Settings.SingleStep = false
			}

			updateAcc -= expectedFrameTime
		}

		//draw debug data
		Draw(scene)

		glfw.SwapBuffers()

		//if one second has passed update the fps and reset the framecount
		if acc > 1 {
			fps = frameCount
			frameCount = 0
			if printFPS {
				fmt.Printf("---\n")
				fmt.Printf("FPS: %v\n", fps)
				fmt.Printf("Update FPS: %v\n", updateFps)
				fmt.Printf("Average frametime: %v\n", acc/float64(fps))
				fmt.Printf("---\n")
			}
			acc -= 1
		}
	}
}
