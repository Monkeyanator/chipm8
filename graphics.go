package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// InitSdlWindow generates window at correct aspect ratio
func InitSdlWindow() (*sdl.Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow(
		"chipm8",
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

	chip.Lock()
	const (
		colorEmpty    = 0x000000
		colorOccupied = 0xFFFFFF
		colorBorder   = 0xD3D3D3
		margin        = 1
	)

	surface, _ := window.GetSurface()
	for i := 0; i < displayRows; i++ {
		for j := 0; j < displayColumns; j++ {
			ind := displayColumns*i + j
			x := int32((ind % displayColumns) * cellSize)
			y := int32(i * cellSize)
			if chip.disp[ind] == 0x0 {
				surface.FillRect(&sdl.Rect{X: x, Y: y, W: cellSize, H: cellSize}, colorBorder)
				surface.FillRect(&sdl.Rect{X: x + margin, Y: y + margin, W: cellSize - margin, H: cellSize - margin}, colorEmpty)
			} else {
				surface.FillRect(&sdl.Rect{X: x, Y: y, W: cellSize, H: cellSize}, colorBorder)
				surface.FillRect(&sdl.Rect{X: x + margin, Y: y + margin, W: cellSize - margin, H: cellSize - margin}, colorOccupied)
			}
		}
	}

	chip.Unlock()
	window.UpdateSurface()

}

func (chip *chip8) SetPixel(x, y uint16) {
	ind := x + y*displayColumns
	chip.disp[ind] = chip.disp[ind] ^ 1
	return

}

func (chip *chip8) IsPixelSet(x, y uint16) bool {
	ind := x + y*displayColumns
	return chip.disp[ind] == 0x1
}
