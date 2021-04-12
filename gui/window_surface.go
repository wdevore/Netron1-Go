package gui

import (
	"Netron1-Go/api"
	"Netron1-Go/simulation"
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// WindowSurface is the GUI and shows the plots and graphs.
// It receives commands for graphing and viewing various graphs.
type WindowSurface struct {
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer
	texture  *sdl.Texture

	rasterBuffer api.IRasterBuffer

	// mouse
	mx int32
	my int32

	running bool
	animate bool
	step    bool

	opened bool

	chFromSim chan string
	chToSim   chan string
}

// NewSurfaceBuffer creates a new viewer and initializes it.
func NewSurfaceBuffer() api.ISurface {
	o := new(WindowSurface)
	o.opened = false
	o.animate = true
	o.step = false
	return o
}

func (ws *WindowSurface) initialize() {
	var err error

	err = sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		panic(err)
	}

	ws.window, err = sdl.CreateWindow("Soft renderer", windowPosX, windowPosY,
		width, height, sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}

	// Using GetSurface requires using window.UpdateSurface() rather than renderer.Present.
	// ws.surface, err = ws.window.GetSurface()
	// if err != nil {
	// 	panic(err)
	// }
	// ws.renderer, err = sdl.CreateSoftwareRenderer(ws.surface)
	// OR create renderer manually
	ws.renderer, err = sdl.CreateRenderer(ws.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	ws.texture, err = ws.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	ws.rasterBuffer = simulation.NewRasterBuffer(width, height)
	// ws.rasterBuffer.EnableAlphaBlending(true)

}

func (ws *WindowSurface) Raster() api.IRasterBuffer {
	return ws.rasterBuffer
}

// Configure view with draw objects
func (ws *WindowSurface) Configure() {
}

// Open shows the viewer and begins event polling
// (host deuron.IHost)
func (ws *WindowSurface) Open() {
	ws.initialize()

	ws.opened = true
}

// SetFont sets the font based on path and size.
func (ws *WindowSurface) SetFont(fontPath string, size int) error {
	var err error
	return err
}

// filterEvent returns false if it handled the event. Returning false
// prevents the event from being added to the queue.
func (ws *WindowSurface) filterEvent(e sdl.Event, userdata interface{}) bool {
	switch t := e.(type) {
	case *sdl.QuitEvent:
		ws.running = false
		return false // We handled it. Don't allow it to be added to the queue.
	case *sdl.MouseMotionEvent:
		ws.mx = t.X
		ws.my = t.Y
		// fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
		// 	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
		return false // We handled it. Don't allow it to be added to the queue.
		// case *sdl.MouseButtonEvent:
		// 	fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
		// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
		// case *sdl.MouseWheelEvent:
		// 	fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
		// 		t.Timestamp, t.Type, t.Which, t.X, t.Y)
	case *sdl.KeyboardEvent:
		if t.State == sdl.PRESSED {
			switch t.Keysym.Scancode {
			case sdl.SCANCODE_ESCAPE:
				ws.chFromSim <- "Exited"
				ws.running = false
			case sdl.SCANCODE_R:
				ws.chToSim <- "run"
			case sdl.SCANCODE_S:
				ws.chToSim <- "reset"
				ws.step = true
				// case 'o':
				// 	// Stop sim
				// 	// simStatus = "Stopping"
				// case 'p':
				// 	// Pause sim
				// 	// simStatus = "Pausing"
				// case 'e':
				// Step sim
				// simStatus = "Stepping"
			}
		}
		// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
		// 	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
		return false
	}

	return true
}

// Run starts the polling event loop. This must run on
// the main thread.
func (ws *WindowSurface) Run(chToSim, chFromSim chan string) {
	ws.chFromSim = chFromSim
	ws.chToSim = chToSim

	// log.Println("Starting viewer polling")
	ws.running = true
	// var simStatus = ""
	var frameStart time.Time
	// var elapsedTime float64
	var loopTime float64

	sleepDelay := 0.0

	// Get a reference to SDL's internal keyboard state. It is updated
	// during sdl.PollEvent()
	keyState := sdl.GetKeyboardState()

	sdl.SetEventFilterFunc(ws.filterEvent, nil)
	ws.rasterBuffer.SetClearColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})

	ws.rasterBuffer.Clear()

	for ws.running {
		frameStart = time.Now()

		sdl.PumpEvents()

		if keyState[sdl.SCANCODE_Z] != 0 {
		}
		if keyState[sdl.SCANCODE_X] != 0 {
		}

		ws.clearDisplay()

		ws.texture.Update(nil, ws.rasterBuffer.Pixels().Pix, ws.rasterBuffer.Pixels().Stride)
		ws.renderer.Copy(ws.texture, nil, nil)

		// fmt.Printf("<%d, %d>\n", ws.mx, ws.my)

		ws.renderer.Present()

		// time.Sleep(time.Millisecond * 5)
		loopTime = float64(time.Since(frameStart).Nanoseconds() / 1000000.0)
		// elapsedTime = float64(time.Since(frameStart).Seconds())

		sleepDelay = math.Floor(framePeriod - loopTime)
		// fmt.Printf("%3.5f ,%3.5f, %3.5f \n", framePeriod, elapsedTime, sleepDelay)
		if sleepDelay > 0 {
			sdl.Delay(uint32(sleepDelay))
			// elapsedTime = framePeriod
		} else {
			// elapsedTime = loopTime
		}

		// f := fmt.Sprintf("%2.2f", 1.0/elapsedTime*1000.0)
		// fmt.Println(f)
	}

	fmt.Println("Run exiting")
}

// Quit stops the gui from running, effectively shutting it down.
func (ws *WindowSurface) Quit() {
	ws.running = false
}

// Close closes the viewer.
// Be sure to setup a "defer x.Close()"
func (ws *WindowSurface) Close() {
	if !ws.opened {
		return
	}
	var err error

	log.Println("Destroying texture")
	err = ws.texture.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying renderer")
	ws.renderer.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying window")
	err = ws.window.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Quiting SDL")
	sdl.Quit()

	if err != nil {
		log.Fatal(err)
	}
}

func (ws *WindowSurface) clearDisplay() {
	ws.window.UpdateSurface()
}

// SetDrawColor --
func (ws *WindowSurface) SetDrawColor(color sdl.Color) {
}

// SetPixel --
func (ws *WindowSurface) SetPixel(x, y int) {
}
