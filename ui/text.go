package ui

import (
	"strings"

	"github.com/clidockermgr/input"
	"github.com/clidockermgr/util"
	"github.com/eiannone/keyboard"
)

type TextView struct {
	ViewImpl
	text     []string
	xpos     uint8
	ypos     uint8
	maxWidth uint8
}

func TextViewNew(text string) *TextView {
	var textView = TextView{text: strings.Split(text, "\n")}
	textView.Init()

	for r := range textView.text {
		textView.maxWidth = uint8(util.Max(int(textView.maxWidth), len(textView.text[r])))
	}

	return &textView
}

func (t *TextView) Draw() {
	var y uint8 = 0
	for v := 0; v < len(t.text); v++ {

		GotoXY(t.rect.x, t.rect.y+y)

		WriteFill(t.text[v], t.rect.w)

		y++
		if y >= t.rect.h {
			break
		}
	}
}

func (t *TextView) ScrollBack() {
	if t.ypos > 0 {
		t.ypos--
	}
}

func (t *TextView) ScrollFwd() {
	t.ypos++
}

func (t *TextView) ScrollLeft() {
	if t.xpos > 0 {
		t.xpos--
	}
}

func (t *TextView) ScrollRight() {
	t.xpos++
}

func (t *TextView) HandleInput(input input.KeyInput) {

	switch input.GetKey() {
	case keyboard.KeyArrowDown:
		t.ScrollFwd()
	case keyboard.KeyArrowUp:
		t.ScrollBack()
	case keyboard.KeyArrowLeft:
		t.ScrollLeft()
	case keyboard.KeyArrowRight:
		t.ScrollRight()
	default:
		t.ViewImpl.HandleInput(input)
	}

}
