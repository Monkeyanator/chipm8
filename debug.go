package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func (chip *chip8) HandleDebugInput(input string) {

	var cmd string
	var count int
	s := strings.Split(input, "x")
	if len(s) == 1 {
		cmd = s[0]
		count = 1
	} else {
		cmd = s[0]
		count, _ = strconv.Atoi(s[1])
	}

	switch cmd {
	case "c", "continue":
		for i := 0; i < count; i++ {
			chip.EmulateNext()
		}
	case "p", "print":
		for i := 0; i < count; i++ {
			chip.printChipState()
		}
	default:
		op, err := decodeTextOpcode(input)
		if err != nil {
			fmt.Println("Error: invalid debug input command")
			return
		}

		chip.EmulateDecodedInstruction(op)
	}
}

func decodeTextOpcode(input string) (uint16, error) {
	hex, err := strconv.ParseInt(input, 0, 32) // 32 bit to fit w/i 4 hex
	if err != nil {
		return 0, err
	}

	return uint16(hex), nil
}

func (chip *chip8) printChipState() {
	fmt16Bit := "[bin] 0b%.16b [hex] 0x%.4x [dec] %d\n"
	fmt8Bit := "[bin] 0b%.8b [hex] 0x%.2x [dec] %d\n"

	// next instruction to execute
	nextOp := chip.DecodeInstruction(chip.pc)
	color.Magenta(fmt.Sprintf(fmt16Bit, nextOp, nextOp, nextOp))
	fmt.Println()

	color.Green(fmt.Sprintf("PC: "+fmt16Bit, chip.pc, chip.pc, chip.pc))
	color.Green(fmt.Sprintf("SP: "+fmt16Bit, chip.sp, chip.sp, chip.sp))
	color.Green(fmt.Sprintf("DT: "+fmt16Bit, chip.dt, chip.dt, chip.dt))
	color.Green(fmt.Sprintf("ST: "+fmt16Bit, chip.st, chip.st, chip.st))
	color.Green(fmt.Sprintf(" I: "+fmt16Bit, chip.I, chip.I, chip.I))
	fmt.Println()

	for i := 0; i < len(chip.reg); i++ {
		color.Red(fmt.Sprintf("V%x: "+fmt8Bit, i, chip.reg[i], chip.reg[i], chip.reg[i]))
	}

	fmt.Println()
}
