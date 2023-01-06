package emulator

import (
	"fmt"
	"os"
)

type Cpu struct {
	display [64][32]uint8
	memory  [4096]uint8
	V       [16]uint8
	stack   []uint16

	pc     uint16
	index  uint16
	opCode uint16

	delayTimer uint8
	soundTimer uint8

	keyMap    map[int]int
	keyInputs [16]bool
}

func newCpu() *Cpu {
	cpu := &Cpu{pc: 0x200}

	cpu.loadFonts()

	return cpu
}

func (c *Cpu) loadFonts() {
	//Loads in font data from 0x00

	var fontset = []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	for i := 0; i < len(fontset); i++ {
		c.memory[i] = fontset[i]
	}
}

func (c *Cpu) loadRom(filePath string) {
	fileData, readErr := os.ReadFile(filePath)
	if readErr != nil {
		fmt.Println(readErr)
	}

	for i := 0; i < len(fileData); i++ {
		c.memory[0x200+i] = fileData[i]
	}
}

func (c *Cpu) cycle() (string, string, bool) {
	c.opCode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	c.pc += 2

	return c.decodeAndExecute()
}

func (c *Cpu) decodeAndExecute() (string, string, bool) {
	identifier := (c.opCode & 0xF000) >> 12
	addr := c.opCode & 0x0FFF
	kk := uint8(c.opCode & 0x00FF)
	x := uint8(c.opCode & 0x0F00 >> 8)
	y := uint8(c.opCode&0x00F0) >> 4
	n := uint8(c.opCode & 0x000F)

	memoryLocation := fmt.Sprintf("%X", c.pc-2)
	instruction := fmt.Sprintf("ERR: #%X", c.opCode)
	drawBool := false

	switch identifier {
	case 0x0:
		if kk == 0xE0 {
			//c.CLS()
			instruction = "CLS"
			drawBool = true
		} else if kk == 0xEE {

		}

	}
}
