package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// RunInputHandler receives input messages from SDL and pushes the result into chip
func (chip *chip8) HandleInput(event sdl.KeyboardEvent) {
	val := SdlKeyToValue(event.Keysym.Sym)
	if val == 0xFF { // means that we found no match for that keycode
		return
	}

	chip.Lock() // protect this op
	if event.Type == sdl.KEYDOWN {
		chip.keys[val] = true
	} else if event.Type == sdl.KEYUP {
		chip.keys[val] = false
	}
	chip.Unlock()
}

// SdlKeyToValue takes an SDL virtual keycode and maps to the hex
func SdlKeyToValue(key sdl.Keycode) uint8 {
	keyToValueMapping := map[sdl.Keycode]uint8{
		sdl.K_1: 0x1,
		sdl.K_2: 0x2,
		sdl.K_3: 0x3,
		sdl.K_4: 0xC,
		sdl.K_q: 0x4,
		sdl.K_w: 0x5,
		sdl.K_e: 0x6,
		sdl.K_r: 0xD,
		sdl.K_a: 0x7,
		sdl.K_s: 0x8,
		sdl.K_d: 0x9,
		sdl.K_f: 0xE,
		sdl.K_z: 0xA,
		sdl.K_x: 0x0,
		sdl.K_c: 0xB,
		sdl.K_v: 0xF,
	}

	val, found := keyToValueMapping[key]
	if !found {
		return 0xFF
	}

	return val
}
