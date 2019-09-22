package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	colorEmpty    = 0x000000
	colorOccupied = 0xFFFFFF
	colorBorder   = 0xD3D3D3
	margin        = 1
)

// Graphics should be used to conduct screen draws
type Graphics interface {
	Render()
	BindBuffer([]byte)
}

type graphics struct {
	window *sdl.Window
	buffer []byte
}

// NewGraphics takes an SDL window and a reference to vidmem
// and renders the screen buffer onto the window through the Render() interface
func NewGraphics(window *sdl.Window) Graphics {
	return &graphics{
		window: window,
	}
}

// Render draws screen buffer onto SDL window
func (g *graphics) Render() {
	for i := 0; i < displayRows; i++ {
		for j := 0; j < displayColumns; j++ {
			ind := displayColumns*i + j
			x := int32((ind % displayColumns) * cellSize)
			y := int32(i * cellSize)
			empty := g.buffer[ind] == 0x0
			renderTile(g.window, x, y, empty, Tiled)
		}
	}
	g.window.UpdateSurface()
}

// BindBuffer takes a slice to bind as vidmem to this Graphics renderer
func (g *graphics) BindBuffer(buf []byte) {
	g.buffer = buf
}

// initSDLWindow generates window at correct aspect ratio
func initSDLWindow() (*sdl.Window, error) {
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

func renderTile(window *sdl.Window, x, y int32, empty, tiled bool) {
	surface, _ := window.GetSurface()

	var fillColor uint32
	if empty {
		fillColor = colorEmpty
	} else {
		fillColor = colorOccupied
	}

	if Tiled {
		surface.FillRect(
			&sdl.Rect{
				X: x + margin,
				Y: y + margin,
				W: cellSize - margin,
				H: cellSize - margin},
			colorEmpty)
		surface.FillRect(
			&sdl.Rect{
				X: x + margin,
				Y: y + margin,
				W: cellSize - margin,
				H: cellSize - margin},
			fillColor)
	} else {
		// not tiled, use standard rect
		surface.FillRect(
			&sdl.Rect{
				X: x,
				Y: y,
				W: cellSize,
				H: cellSize},
			fillColor)
	}
}
