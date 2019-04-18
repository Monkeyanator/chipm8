package main

import (
	"testing"
)

func generateMockChip8() *chip8 {
	chip := &chip8{}
	chip.InitChip8Registers()
	render := make(chan bool, 1) // unbuffered channels block, for tests use a size 1 bufer
	sound := make(chan bool, 1)
	chip.render = render
	chip.sound = sound
	return chip
}

func TestCLS(t *testing.T) {
	chip := generateMockChip8()

	chip.disp[0] = 0xF1
	chip.disp[1] = 0x5A
	chip.disp[1023] = 0xFF

	chip.EmulateDecodedInstruction(0x00E0)
	for i := 0; i < len(chip.disp); i++ {
		if chip.disp[i] != 0x0 {
			t.Error("Test Failed: expected 0x00E0 instruction to clear disp bits, but did not")
		}
	}

}

func TestRet(t *testing.T) {

	chip := generateMockChip8()
	chip.pc = 0x8
	chip.sp = 2
	chip.stack[0] = 0x1
	chip.stack[1] = 0x2
	chip.stack[2] = 0x3

	chip.EmulateDecodedInstruction(0x00EE)

	if chip.pc != 0x5 {
		t.Errorf("Test Failed: expected program counter to update, found value 0x%x", chip.pc)
	}

	if chip.sp != 1 {
		t.Errorf("Test Failed: expected stack pointer to decrement, found 0x%x", chip.sp)
	}
}

func TestSetPC(t *testing.T) {
	var tests = []struct {
		opcode uint16
		result address
	}{
		{0x1015, 0x015},
		{0x1000, 0x000},
		{0x1FFF, 0xFFF},
	}

	for _, test := range tests {
		chip := generateMockChip8()
		chip.EmulateDecodedInstruction(test.opcode)
		if chip.pc != test.result {
			t.Errorf("Test Failed: expected program counter to update to 0x%x, found instead 0x%x", test.result, chip.pc)
		}
	}
}

func TestCall(t *testing.T) {
	type setup struct {
		stack []address
		sp    uint8
		pc    address
	}

	type result struct {
		sp       uint8
		pc       address
		stackTop address
	}

	var tests = []struct {
		setup setup
		res   result
	}{
		{setup{[]address{0x1, 0x2, 0x3}, 2, 0x200}, result{3, 0x00F, 0x200}},
	}

	for _, test := range tests {
		chip := generateMockChip8()
		chip.sp = test.setup.sp
		chip.pc = test.setup.pc
		for i, s := range test.setup.stack {
			chip.stack[i] = s
		}
		chip.EmulateDecodedInstruction(0x200F)

		if chip.pc != test.res.pc {
			t.Errorf("Expected program counter to update to 0x%x, got 0x%x", test.res.pc, chip.pc)
		}

		if chip.sp != test.res.sp {
			t.Errorf("Expected stack pointer to update to %d, got %d", test.res.sp, chip.sp)
		}

		actualStackTop := chip.stack[chip.sp]
		if actualStackTop != test.res.stackTop {
			t.Errorf("Expected address at top of stack to be 0x%x, got 0x%x", test.res.stackTop, chip.stack[chip.sp])
		}

	}

}
