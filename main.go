package main

import "fmt"

var opCodes uint16
var memory [4096]uint8
var V [16]uint8 //Register V0-VF

func main() {
	var a uint8 = 0xA2
	b := 0xFA

	c := uint16(a)<<8 | uint16(b)
	d := c & 0xF000

	fmt.Println(d >> 12)
}
