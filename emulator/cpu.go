package emulator

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
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
			c.CLS()
			instruction = "CLS"
			drawBool = true
		} else if kk == 0xEE {
			c.RET()
			instruction = "RET"
		}

	case 0x1:
		c.JP(addr)
		instruction = fmt.Sprintf("JP #%X", addr)

	case 0x2:
		c.CALL(addr)
		instruction = fmt.Sprintf("CALL #%X", addr)

	case 0x3:
		c.SEVx(x, kk)
		instruction = fmt.Sprintf("SE V%X #%X", x, kk)
	case 0x4:
		c.SNEVx(x, kk)
		instruction = fmt.Sprintf("SNE V%X #%X", x, kk)
	case 0x5:
		c.SEVxVy(x, y)
		instruction = fmt.Sprintf("SE V%X V%X", x, y)
	case 0x6:
		c.LDVx(x, kk)
		instruction = fmt.Sprintf("LD V%X #%X", x, kk)
	case 0x7:
		c.ADDVx(x, kk)
		instruction = fmt.Sprintf("ADD V%X #%X", x, kk)
	case 0x8:
		switch n {
		case 0x0:
			c.LDVxVy(x, y)
			instruction = fmt.Sprintf("LD V%X V%X", x, y)
		case 0x1:
			c.ORVxVy(x, y)
			instruction = fmt.Sprintf("OR V%X V%X", x, y)
		case 0x2:
			c.ANDVxVy(x, y)
			instruction = fmt.Sprintf("AND V%X V%X", x, y)
		case 0x3:
			c.XORVxVy(x, y)
			instruction = fmt.Sprintf("XOR V%X V%X", x, y)
		case 0x4:
			c.ADDVxVy(x, y)
			instruction = fmt.Sprintf("ADD V%X V%X", x, y)
		case 0x5:
			c.SUBVxVy(x, y)
			instruction = fmt.Sprintf("SUB V%X V%X", x, y)
		case 0x6:
			c.SHRVx(x)
			instruction = fmt.Sprintf("SHR V%X", x)
		case 0x7:
			c.SUBNVxVy(x, y)
			instruction = fmt.Sprintf("SUBN V%X V%X", x, y)
		case 0xE:
			c.SHLVx(x)
			instruction = fmt.Sprintf("SHL V%X", x)
		}
	case 0x9:
		c.SNEVxVy(x, y)
		instruction = fmt.Sprintf("SNE V%X V%X", x, y)
	case 0xA:
		c.LDI(addr)
		instruction = fmt.Sprintf("LD I #%X", addr)
	case 0xB:
		c.JPV(addr)
		instruction = fmt.Sprintf("JP V0 #%X", addr)
	case 0xC:
		c.RNDVx(x, kk)
		instruction = fmt.Sprintf("RND V%X #%X", x, kk)
	case 0xD:
		instruction = fmt.Sprintf("DRW V%X V%X #%X", x, y, n)
		c.DRW(x, y, n)
		drawBool = true
	case 0xE:
		if n == 0xE {
			c.SKPVx(x)
			instruction = fmt.Sprintf("SKP V%X", x)
		} else if n == 0x1 {
			c.SKNPVx(x)
			instruction = fmt.Sprintf("SKNP V%X", x)
		}
	case 0xF:
		switch kk {
		case 0x07:
			c.LDVxDT(x)
			instruction = fmt.Sprintf("LD V%X DT", x)
		case 0x0A:
			c.LDVxK(x)
			instruction = fmt.Sprintf("LD V%X K", x)
		case 0x15:
			c.LDDTVx(x)
			instruction = fmt.Sprintf("LD DT V%X", x)
		case 0x18:
			c.LDSTVx(x)
			instruction = fmt.Sprintf("LD ST V%X", x)
		case 0x1E:
			c.ADDIVx(x)
			instruction = fmt.Sprintf("ADD I V%X", x)
		case 0x29:
			c.LDFVx(x)
			instruction = fmt.Sprintf("LD F V%X", x)
		case 0x33:
			c.LDBVx(x)
			instruction = fmt.Sprintf("LD B V%X", x)
		case 0x55:
			c.LDIVx(x)
			instruction = fmt.Sprintf("LD I V%X", x)
		case 0x65:
			c.LDVxI(x)
			instruction = fmt.Sprintf("LD V%X I", x)
		}
	}

	return memoryLocation, instruction, drawBool
}

func (c *Cpu) handleKeypress(scancode sdl.Scancode, keystate bool) {
	//Use the keymap to correctly handle keydown and keyups
	c.keyInputs[c.keyMap[int(scancode)]] = keystate
}

func (c *Cpu) LDVxI(x uint8) {

}

func (c *Cpu) LDIVx(x uint8) {

}

func (c *Cpu) LDBVx(x uint8) {

}

func (c *Cpu) LDFVx(x uint8) {

}

func (c *Cpu) CLS() {

}

func (c *Cpu) RET() {

}

func (c *Cpu) JP(addr uint16) {

}

func (c *Cpu) CALL(addr uint16) {

}

func (c *Cpu) SEVx(x uint8, kk uint8) {

}

func (c *Cpu) SNEVx(x uint8, kk uint8) {

}

func (c *Cpu) SEVxVy(x uint8, y uint8) {

}

func (c *Cpu) LDVx(x uint8, kk uint8) {

}

func (c *Cpu) ADDVx(x uint8, kk uint8) {

}

func (c *Cpu) LDVxVy(x uint8, y uint8) {

}

func (c *Cpu) ORVxVy(x uint8, y uint8) {

}

func (c *Cpu) ANDVxVy(x uint8, y uint8) {

}

func (c *Cpu) XORVxVy(x uint8, y uint8) {

}

func (c *Cpu) ADDVxVy(x uint8, y uint8) {

}

func (c *Cpu) SUBVxVy(x uint8, y uint8) {

}

func (c *Cpu) SHRVx(x uint8) {

}

func (c *Cpu) SUBNVxVy(x uint8, y uint8) {

}

func (c *Cpu) SHLVx(x uint8) {

}

func (c *Cpu) SNEVxVy(x uint8, y uint8) {

}

func (c *Cpu) LDI(addr uint16) {

}

func (c *Cpu) JPV(addr uint16) {

}

func (c *Cpu) RNDVx(x uint8, kk uint8) {

}

func (c *Cpu) DRW(x uint8, y uint8, n uint8) {

}

func (c *Cpu) SKPVx(x uint8) {

}

func (c *Cpu) SKNPVx(x uint8) {

}

func (c *Cpu) LDVxDT(x uint8) {

}

func (c *Cpu) LDVxK(x uint8) {

}

func (c *Cpu) LDDTVx(x uint8) {

}

func (c *Cpu) LDSTVx(x uint8) {

}

func (c *Cpu) ADDIVx(x uint8) {

}
