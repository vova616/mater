package main

import (
	"flag"
	"fmt"
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/jteeuwen/glfw"
	"github.com/teomat/mater/collision"
	"log"
	//importing so the components can register themselves
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

var console Console

var space *collision.Space

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
		//create empty space
		space = collision.NewSpace()
	}

	//reload settings so they take effect
	reloadSettings()

	//set callbacks
	{
		glfw.SetWindowSizeCallback(OnResize)
		glfw.SetKeyCallback(OnKey)
	}

	//init debug console
	console.Init()

	//load savefile passed from the commandline if any
	if flags.file != "" {
		err := loadSpace(flags.file)
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
				space.Step(expectedFrameTime)
				Settings.SingleStep = false
			}

			updateAcc -= expectedFrameTime
		}

		//draw debug data
		Draw()

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
