package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// InitSdlWindow generates window at correct aspect ratio
func InitSdlWindow() (*sdl.Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"Chip8",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		displayColumns*cellSize,
		displayRows*cellSize,
		sdl.WINDOW_SHOWN)

	if err != nil {
		return nil, err
	}

	return window, nil
}

// RenderChip8 takes an SDL window and chip8 state, and renders
func RenderChip8(window *sdl.Window, chip *chip8) {

	const (
		colorEmpty    = 0x000000
		colorOccupied = 0xFFFFFF
		cellSize      = 10
	)

	surface, _ := window.GetSurface()
	for i := 0; i < displayRows; i++ {
		for j := 0; j < displayColumns; j++ {
			ind := displayColumns*i + j
			x := int32((ind % displayColumns) * cellSize)
			y := int32(i * cellSize)
			if chip.disp[ind] == 0x0 {
				surface.FillRect(&sdl.Rect{X: x, Y: y, W: cellSize, H: cellSize}, colorEmpty)
			} else {
				surface.FillRect(&sdl.Rect{X: x, Y: y, W: cellSize, H: cellSize}, colorOccupied)
			}
		}
	}

	window.UpdateSurface()

}
