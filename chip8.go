package main

import (
	"bufio"
	"os"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	displayRows    = 32
	displayColumns = 64
	memSize        = 4096
	stackDepth     = 16
	programOffset  = 0x200
	hz             = time.Duration(600) // m/s b/w emulation steps
	cellSize       = 15
)

type address uint16

type chip8 struct {
	sync.Mutex
	reg    [16]byte // registers, reg[15] is carry
	stack  [stackDepth]address
	pc     address                            // program counter
	sp     uint8                              // stack pointer
	dt     byte                               // delay timer
	st     byte                               // sound timer
	I      uint16                             // special reg, stores addresses
	mem    [memSize]byte                      // rom and program work ram
	disp   [displayRows * displayColumns]byte // graphics mem
	keys   [16]bool                           // stores keypress state
	input  chan sdl.KeyboardEvent
	render chan bool
	sound  chan bool
}

func (chip *chip8) Init() {
	chip.InitChip8Registers()
	chip.InitCharset()
}

func (chip *chip8) InitChip8Registers() {
	// note that reg defaults elems to 0x00
	chip.dt = 0x0
	chip.st = 0x0
	chip.pc = programOffset // program execution starts here by convention
}

func (chip *chip8) LoadProgram(path string) {
	progBuff, err := createProgramBuffer(path)
	if err != nil {
		panic(progBuff)
	}

	// load the mem in at offset
	for i := 0; i < len(progBuff); i++ {
		chip.mem[programOffset+i] = progBuff[i]
	}
}

/* register manipulation */
func (chip *chip8) SetRegister(index uint8, value byte) {
	chip.reg[index] = value
}

func (chip *chip8) ReadRegister(index uint8) byte {
	return chip.reg[index]
}

/* pc manipulation */
func (chip *chip8) IncrementPC() {
	chip.pc += 2
}

func (chip *chip8) SetPC(val address) {
	chip.pc = address(val)
}

/* helpers */
func createProgramBuffer(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	stats, statsErr := f.Stat()
	if statsErr != nil {
		return nil, statsErr
	}
	size := stats.Size()
	bytes := make([]byte, size)

	buff := bufio.NewReader(f)
	_, err = buff.Read(bytes)

	return bytes, err
}
