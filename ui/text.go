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

	var firstLine = t.ypos
	var lastLine = t.ypos + t.rect.h - 1

	var length = uint8(len(t.text))

	if lastLine >= length {
		lastLine = length - 1
		firstLine = lastLine - t.rect.h
	}

	for v := firstLine; v <= lastLine; v++ {

		GotoXY(t.rect.x, t.rect.y+y)

		var line = t.text[v]
		if t.xpos < uint8(len(line)) {
			line = line[t.xpos:]
		} else {
			line = ""
		}

		WriteFill(line, t.rect.w)

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
	if t.ypos+t.rect.h < uint8(len(t.text))-1 {
		t.ypos++
	}
}

func (t *TextView) ScrollLeft() {
	if t.xpos > 0 {
		t.xpos--
	}
}

func (t *TextView) ScrollRight() {
	if t.xpos+t.rect.w < t.maxWidth {
		t.xpos++
	}
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
