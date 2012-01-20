package main

import (
	"flag"
	"fmt"
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/jteeuwen/glfw"
	. "github.com/teomat/mater"
	"github.com/teomat/mater/vect"
	"log"
	//importing so the components can register themselves
	_ "github.com/teomat/mater/components"

	"os"
	"runtime/pprof"
	"github.com/teomat/mater/camera"
)

var flags = struct {
	startPaused   bool
	help          bool
	dbg           bool
	file          string
	buildExamples bool
}{}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func init() {
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

var MainCamera *camera.Camera
var console Console

func main() {
	log.SetFlags(log.Lshortfile)

	//parse flags
	flag.Parse()
	if flags.help {
		flag.PrintDefaults()
		return
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	loadSettingsFile()

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

	//setup default camera
	{
		MainCamera = new(camera.Camera)
		cam := MainCamera
		camera.ScreenSize = vect.Vect{float64(wx), float64(wy)}
		cam.Transform.Position = vect.Vect{0, 0}
		cam.Scale = vect.Vect{32, 32}
		cam.Transform.SetAngle(0)
	}

	mater := new(Mater)
	mater.Init()

	reloadSettings(mater)

	Settings.Paused = flags.startPaused

	mater.Callbacks.OnNewComponent = OnNewComponent

	console.Init(mater)

	if flags.file != "" {
		err := mater.LoadScene(flags.file)
		Settings.Paused = true
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

	//glfw config
	{
		glfw.SetSwapInterval(1)
		glfw.SetWindowTitle("mater test")
		glfw.SetWindowSizeCallback(func(w, h int) {OnResize(w, h)})
		glfw.SetKeyCallback(func(k, s int) { OnKey(mater, k, s) })
	}

	//init opengl
	{
		gl.ClearColor(0, 0, 0, 0)
		gl.Enable(gl.BLEND)
		//gl.BlendFunc (gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
		gl.Enable(gl.TEXTURE_2D)
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
	Settings.Running = true

	for Settings.Running && glfw.WindowParam(glfw.Opened) == 1 {
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
		case command := <-console.Command:
			console.ExecuteCommand(command)
		default:
		}

		for updateAcc >= expectedFrameTime {
			updateFrameCount++

			if acc > 1 {
				updateFps = updateFrameCount
				updateFrameCount = 0
			}

			if !Settings.Paused || Settings.SingleStep {
				mater.Update(expectedFrameTime)
				if Settings.SingleStep {
					Settings.SingleStep = false
				}
			}

			updateAcc -= expectedFrameTime
		}

		Draw(mater)

		glfw.SwapBuffers()

		if acc > 1 {
			fps = frameCount
			frameCount = 0
			if !Settings.Paused && printFPS {
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
