package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	var debugMode bool
	var programPath string
	flag.BoolVar(&debugMode, "debug", false, "Program enters into interactive debugger mode, wherein user can step through opcodes one at a time")
	flag.StringVar(&programPath, "prog", "", "Path to the chip8 machine code to emulate")
	flag.Parse()

	window, err := InitSdlWindow()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer sdl.Quit()

	chip := &chip8{}
	chip.Init()
	chip.LoadProgram(programPath) // should return err

	if debugMode {
		debugLoop(window, chip)
	} else {
		mainLoop(window, chip)
	}

}

func mainLoop(window *sdl.Window, chip *chip8) {

	input := make(chan sdl.KeyboardEvent)
	render := make(chan bool)
	chip.input = input
	chip.render = render

	done := make(chan bool) // should not be here
	go RunInputHandler(chip, input, done)

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	go func() {
		for {
			<-render
			RenderChip8(window, chip)
		}
	}()

	go func() {
		for {
			chip.EmulateNext()
			time.Sleep(time.Millisecond * speed)
		}
	}()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				// keyinput <- event.(*sdl.KeyboardEvent)
				input <- *event.(*sdl.KeyboardEvent)
			}
		}
	}

}

func debugLoop(window *sdl.Window, chip *chip8) {

	input := make(chan sdl.KeyboardEvent)
	chip.input = input

	done := make(chan bool) // should not be here
	go RunInputHandler(chip, input, done)

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

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				// keyinput <- event.(*sdl.KeyboardEvent)
				input <- *event.(*sdl.KeyboardEvent)
			}
		}
	}

}
