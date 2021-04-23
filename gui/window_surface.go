package gui

import (
	"Netron1-Go/api"
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
	model        api.IModel

	// mouse
	mx int32
	my int32

	running bool
	animate bool
	step    bool

	opened bool
	ready  bool

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

func (ws *WindowSurface) Raster() api.IRasterBuffer {
	return ws.rasterBuffer
}

// Open shows the viewer and begins event polling
// (host deuron.IHost)
func (ws *WindowSurface) Open(model api.IModel) {
	var err error
	ws.model = model

	err = sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		panic(err)
	}

	mp := model.Properties()

	ws.window, err = sdl.CreateWindow("Soft renderer",
		int32(mp.WindowPosX()), int32(mp.WindowPosY()),
		int32(mp.Width()), int32(mp.Height()), sdl.WINDOW_SHOWN)

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

	w := int32(mp.Width())
	h := int32(mp.Height())

	ws.texture, err = ws.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, w, h)
	if err != nil {
		panic(err)
	}

	ws.rasterBuffer = NewRasterBuffer(int(w), int(h))
	ws.rasterBuffer.EnableAlphaBlending(true)

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
			case sdl.SCANCODE_E:
				ws.chToSim <- "step"
			case sdl.SCANCODE_P:
				ws.chToSim <- "pause"
			case sdl.SCANCODE_U:
				ws.chToSim <- "resume"
			case sdl.SCANCODE_T:
				ws.chToSim <- "stop"
			case sdl.SCANCODE_A:
				ws.chToSim <- "status"
			case sdl.SCANCODE_S:
				ws.chToSim <- "reset"
				ws.step = true
			case sdl.SCANCODE_K: // decrease acceptible rate
				ws.chToSim <- "inc accept"
			case sdl.SCANCODE_L: // increase accetable rate
				ws.chToSim <- "dec accept"
			case sdl.SCANCODE_N: // decrease drop rate
				ws.chToSim <- "dec drop"
			case sdl.SCANCODE_M: // increase drop rate
				ws.chToSim <- "inc drop"
			case sdl.SCANCODE_COMMA: // decrease step size
				ws.chToSim <- "dec size"
			case sdl.SCANCODE_PERIOD: // increase step size
				ws.chToSim <- "inc size"
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

	ws.running = true
	var frameStart time.Time
	var loopTime float64

	sleepDelay := 0.0

	// Get a reference to SDL's internal keyboard state. It is updated
	// during sdl.PollEvent()
	keyState := sdl.GetKeyboardState()

	sdl.SetEventFilterFunc(ws.filterEvent, nil)
	ws.rasterBuffer.SetClearColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})

	ws.rasterBuffer.Clear()

	mp := ws.model.Properties()
	// leftT := sdl.Rect{X: 0, Y: 0, W: int32(mp.Width()) / 2, H: int32(mp.Height())}
	leftT := sdl.Rect{X: 0, Y: 0, W: int32(mp.Width() * mp.Scale()), H: int32(mp.Height() * mp.Scale())}
	// rightT := sdl.Rect{X: int32(mp.Width() / 2), Y: 0, W: int32(mp.Width() / 2), H: int32(mp.Height())}

	for ws.running {
		frameStart = time.Now()

		sdl.PumpEvents()

		if keyState[sdl.SCANCODE_Z] != 0 {
		}
		if keyState[sdl.SCANCODE_X] != 0 {
		}

		ws.clearDisplay()

		pixs := ws.rasterBuffer.BackPixels()
		ws.texture.Update(nil, pixs.Pix, pixs.Stride)
		// ws.renderer.Copy(ws.texture, nil, nil)
		ws.renderer.Copy(ws.texture, &leftT, &leftT)
		// ws.renderer.Copy(ws.texture, &rightT, &rightT)

		// fmt.Printf("<%d, %d>\n", ws.mx, ws.my)

		ws.renderer.Present()

		loopTime = float64(time.Since(frameStart).Nanoseconds() / 1000000.0)

		sleepDelay = math.Floor(framePeriod - loopTime)
		if sleepDelay > 0 {
			sdl.Delay(uint32(sleepDelay))
		} else {
		}
	}

	fmt.Println("Run exiting")
}

// Update
func (ws *WindowSurface) Update(state bool) {
	ws.rasterBuffer.Swap()
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
