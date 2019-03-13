package main

import (
	"testing"

	"github.com/veandco/go-sdl2/sdl"
)

func TestSdlKeyToValue(t *testing.T) {
	var tests = []struct {
		keycode sdl.Keycode
		result  uint8
	}{
		{sdl.K_2, 0x2},
		{sdl.K_q, 0x4},
		{sdl.K_BACKSPACE, 0xFF},
	}

	for _, test := range tests {
		val := SdlKeyToValue(test.keycode)
		if val != test.result {
			t.Errorf("Expected %d, got %d", test.result, val)
		}
	}
}
