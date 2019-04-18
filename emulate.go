package main

import (
	"math/rand"
	"time"
)

func (chip *chip8) EmulateNext() {
	decodedInstruction := chip.DecodeInstruction(chip.pc)
	chip.EmulateDecodedInstruction(decodedInstruction)
}

func (chip *chip8) EmulateDecodedInstruction(op uint16) {

	// flag showing whether we should inc program counter
	retflag := true

	// mask to get first hex digit
	switch op & 0xF000 {

	case 0x0000:

		switch op & 0xFF {

		case 0xE0: // CLS
			for i := 0; i < len(chip.disp); i++ {
				chip.disp[i] = 0x00
			}

			chip.render <- true // seems wrong

		case 0xEE: // RET
			chip.SetPC(chip.stack[chip.sp])
			chip.sp--
		}

	case 0x1000: // SET PC
		chip.SetPC(address(op & 0x0FFF))
		retflag = false

	case 0x2000: // CALL --> 2NNN
		chip.sp++
		chip.stack[chip.sp] = chip.pc // must return past the last instruction
		chip.SetPC(address(op & 0x0FFF))
		retflag = false

	case 0x3000: // SE
		targetReg := uint8((op >> 8) & 0x000F)
		compareVal := uint8(op & 0x00FF)
		if chip.ReadRegister(targetReg) == compareVal {
			chip.IncrementPC()
		}

	case 0x4000: // SNE
		targetReg := uint8((op >> 8) & 0x000F)
		compareVal := uint8(op & 0x00FF)
		if chip.ReadRegister(targetReg) != compareVal {
			chip.IncrementPC()
		}

	case 0x5000: // SE Vx Vy -> 5XY0
		rxval := chip.ReadRegister(uint8((op >> 8) & 0x000F))
		ryval := chip.ReadRegister(uint8((op >> 4) & 0x000F))
		if rxval == ryval {
			chip.IncrementPC()
		}

	case 0x6000: // LD Vx, byte --> 6xKK
		targetReg := uint8((op >> 8) & 0x000F)
		ldVal := uint8(op & 0x00FF)
		chip.SetRegister(targetReg, ldVal)

	case 0x7000: // ADD Vx, byte --> 7xKK
		targetReg := uint8((op >> 8) & 0x000F)
		kk := uint8(op & 0x00FF)
		sum := chip.ReadRegister(targetReg) + kk
		chip.SetRegister(targetReg, sum)

	case 0x8000: // these take form 8xyZ
		rx := uint8((op >> 8) & 0x000F)
		ry := uint8((op >> 4) & 0x000F)
		rxval := chip.ReadRegister(rx)
		ryval := chip.ReadRegister(ry)

		switch op & 0x0F {

		case 0x0:
			chip.SetRegister(rx, ryval)

		case 0x1: // OR Vx, Vy
			chip.SetRegister(rx, rxval|ryval)

		case 0x2: // AND Vx, Vy
			chip.SetRegister(rx, rxval&ryval)

		case 0x3: // XOR Vx, Vy
			chip.SetRegister(rx, rxval^ryval)

		case 0x4: // ADD Vx, Vy
			sum := uint16(rxval) + uint16(ryval)
			trimmedSum := byte(sum)
			chip.SetRegister(rx, trimmedSum)
			chip.setRegisterOnCondition(0xf, sum > 255)

		case 0x5: // SUB Vx, Vy
			chip.setRegisterOnCondition(0xf, rxval > ryval)
			chip.SetRegister(rx, rxval-ryval)

		case 0x6: // SHR Vx {, Vy}
			chip.setRegisterOnCondition(0xf, rxval&0x1 == 0x1)
			chip.SetRegister(rx, rxval/2)

		case 0x7: // SUBN Vx, Vy
			chip.setRegisterOnCondition(0xf, ryval > rxval)
			chip.SetRegister(rx, ryval-rxval)

		case 0xE: // SHL Vx {, Vy}
			chip.setRegisterOnCondition(0xf, (rxval>>7)&0x1 == 0x1)
			chip.SetRegister(rx, rxval*2)
		}

	case 0x9000: // SNE Vx, Vy --> 9xy0
		rxval := chip.ReadRegister(uint8((op >> 8) & 0x000F))
		ryval := chip.ReadRegister(uint8((op >> 4) & 0x000F))
		if rxval != ryval {
			chip.IncrementPC()
		}

	case 0xA000: // LD I, addr --> ANNN
		chip.I = 0x0FFF & op

	case 0xB000: // JP V0, addr --> BNNN
		baseAddr := address(op & 0x0FFF)
		valV0 := chip.ReadRegister(0x0)
		chip.SetPC(baseAddr + address(valV0)) // can we make this conversion?
		retflag = false                       // noinc

	case 0xC000: // RND Vx, byte --> CxKK
		targetReg := uint8((op >> 8) & 0x000F)
		val := uint8(op & 0x00FF)
		rand.Seed(time.Now().UnixNano())
		result := uint8(rand.Intn(256)) & val
		chip.SetRegister(targetReg, byte(result))

	case 0xD000: // DRW Vx, Vy --> 0xDXYN
		rxval := uint16(chip.ReadRegister(uint8((op >> 8) & 0x000F)))
		ryval := uint16(chip.ReadRegister(uint8((op >> 4) & 0x000F)))
		n := uint16(op & 0x000F)

		var b byte
		chip.SetRegister(0xf, 0x0)
		for yline := uint16(0); yline < n; yline++ {
			b = chip.mem[chip.I+uint16(yline)]
			for xline := uint16(0); xline < 8; xline++ {

				xWrapped := (xline + rxval) % displayColumns
				yWrapped := (yline + ryval) % displayRows

				if b&(0x80>>xline) != 0 {
					if chip.IsPixelSet(xWrapped, yWrapped) {
						chip.SetRegister(0xf, 0x1)
					}
					chip.SetPixel(xWrapped, yWrapped)
				}

			}
		}

		chip.render <- true

	case 0xE000: // SKP Vx --> Ex9E
		targetReg := uint8((op >> 8) & 0x000F)

		switch op & 0xFF {

		case 0x9E:
			if chip.keys[chip.ReadRegister(targetReg)] {
				chip.IncrementPC()
			}

		case 0xA1:
			if !chip.keys[chip.ReadRegister(targetReg)] {
				chip.IncrementPC()
			}
		}

	case 0xF000:
		targetReg := uint8((op >> 8) & 0x000F)

		// swith over last two digits of opcode for 0xF___
		switch op & 0xFF {

		case 0x07: // LD Vx, DT --> FX07
			chip.SetRegister(targetReg, chip.dt)

		case 0x0A: // LD Vx, K --> FX0A
			// there should be input channel that some keypress handler writes to,
			// and this should wait on that channel, for now scanf
			event := <-chip.input
			value := SdlKeyToValue(event.Keysym.Sym)
			chip.SetRegister(targetReg, value)

		case 0x15: // LD DT, Vx --> FX15
			chip.dt = chip.ReadRegister(targetReg)

		case 0x18:
			chip.st = chip.ReadRegister(targetReg)

		case 0x1E:
			chip.I = uint16(chip.ReadRegister(targetReg)) + chip.I

		case 0x29:
			hexDigToPointTo := chip.ReadRegister(targetReg)
			chip.I = uint16(hexDigToPointTo) * 0x5

		case 0x33: // LD B, Vx --> FX33 (BCD into I, I+1, I+2)
			val := uint8(chip.ReadRegister(targetReg))
			chip.mem[chip.I] = byte(val / 100)
			chip.mem[chip.I+1] = byte((val / 10) % 10)
			chip.mem[chip.I+2] = byte(val % 10)

		case 0x55: // LD [I], Vx
			for i := 0; i <= int(targetReg); i++ {
				chip.mem[int(chip.I)+i] = chip.ReadRegister(byte(i))
			}

		case 0x65: // LD Vx, [I]
			for i := 0; i <= int(targetReg); i++ {
				chip.SetRegister(uint8(i), chip.mem[int(chip.I)+i])
			}

		}

	}

	// increment program counter if needed, and decrement
	if retflag {
		chip.IncrementPC()
	}

	// update timers
	if chip.st > 0 {
		chip.sound <- true
		chip.st = chip.st - 1
	} else {
		chip.sound <- false
	}

	if chip.dt > 0 {
		chip.dt = chip.dt - 1
	}

}

// DecodeInstruction takes an index and decodes the contained bytes as opcode
func (chip *chip8) DecodeInstruction(ind address) uint16 {
	decodedInstruction := (uint16(chip.mem[ind]) << 8) | uint16(chip.mem[ind+1])
	return decodedInstruction
}

func (chip *chip8) setRegisterOnCondition(reg uint8, predicate bool) {
	if predicate {
		chip.SetRegister(reg, 1)
	} else {
		chip.SetRegister(reg, 0)
	}
}

// Max provides basic byte min
func Max(a, b byte) byte {
	if a > b {
		return a
	}
	return b
}
