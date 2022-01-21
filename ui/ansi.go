package ui

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func ScreenSize() (uint8, uint8) {
	// https://stackoverflow.com/questions/16569433/get-terminal-size-in-go
	ws := &winsize{}

	retCode, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(err)
	}

	return uint8(ws.Col), uint8(ws.Row)
}

func GotoXY(x, y uint8) {
	fmt.Printf("\u001b[%d;%dH", y, x)
}

func Bold() {
	fmt.Print("\u001b[1m")
}

func Foreground(color uint8) {
	fmt.Printf("\u001b[38;5;%dm", color)
}

func Background(color uint8) {
	fmt.Printf("\u001b[48;5;%dm", color)
}

func WriteFill(text string, length uint8) {
	if len(text) > int(length) {
		fmt.Print(text[0:length])
	} else {
		fmt.Printf("%s%s", text, strings.Repeat(" ", int(length)-len(text)))
	}
}

func WriteV(char string, x, y, length uint8) {
	for i := 0; i < int(length); i++ {
		GotoXY(x, y+uint8(i))
		fmt.Print(char)
	}
}

func ClearScreen() {
	fmt.Print("\u001b[2J")
}

func CursorOff() {
	fmt.Print("\033[?25l")
}

func CursorOn() {
	fmt.Print("\033[?25h")
}

func UnderlineOn() {
	fmt.Print("\u001b[4m")
}

func Reset() {
	fmt.Print("\u001b[0m")
}
