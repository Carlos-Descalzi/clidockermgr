package ui

import (
	"strings"

	"github.com/eiannone/keyboard"
)

type TextView struct {
	ViewImpl
	text []string
	xpos int
	ypos int
}

func TextViewNew(text string) *TextView {
	return &TextView{text: strings.Split(text, "\n")}
}

func (t *TextView) Draw() {
}

func (t *TextView) HandleInput(key keyboard.Key) {

	switch key {
	case keyboard.KeyArrowDown:
		t.ypos++
	case keyboard.KeyArrowUp:
		if t.ypos > 0 {
			t.ypos--
		}
	default:
		t.ViewImpl.HandleInput(key)
	}

}
