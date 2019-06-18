package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var DebugMode bool
var ProgramPath string
var Tiled bool

func main() {

	flag.BoolVar(&DebugMode, "debug", false, "Program enters into interactive debugger mode, wherein user can step through opcodes one at a time")
	flag.StringVar(&ProgramPath, "prog", "", "Path to the chip8 machine code to emulate")
	flag.BoolVar(&Tiled, "tiled", false, "Determines whether the chip8 display will include a border around pixels")
	flag.Parse()

	window, err := InitSdlWindow()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer sdl.Quit()

	chip := &chip8{}
	chip.Init()
	chip.LoadProgram(ProgramPath) // should return err

	if DebugMode {
		debugLoop(window, chip)
	} else {
		mainLoop(window, chip)
	}

}

func mainLoop(window *sdl.Window, chip *chip8) {

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	// controls emulation ticks
	go func() {
		ticker := time.NewTicker(time.Second / hz)
		for {
			<-ticker.C
			chip.tick <- true
		}
	}()

	// this handles main emulation logic
	go chip.emulationLoop(window)
	os.Exit(chip.SDLLoop())
}

func debugLoop(window *sdl.Window, chip *chip8) {

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for {
			fmt.Print("[CH-I-P8]> ")
			scanner.Scan()
			input := scanner.Text()
			chip.HandleDebugInput(input)
			RenderChip8(window, chip)
		}
	}()

	// this handles main emulation logic
	go chip.emulationLoop(window)
	os.Exit(chip.SDLLoop())
}

func (chip *chip8) emulationLoop(window *sdl.Window) {
	for {
		select {
		case <-chip.render:
			RenderChip8(window, chip)

		case <-chip.sound:
			break

		case input := <-chip.input:
			chip.HandleInput(input)

		case <-chip.tick:
			chip.EmulateNext()
		}
	}
}

func (chip *chip8) SDLLoop() int {
	// SDL poll loop
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				chip.input <- *event.(*sdl.KeyboardEvent)
			}
		}
	}
	return 0
}
