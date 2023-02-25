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
	xpos     uint16
	ypos     uint16
	maxWidth uint16
}

func TextViewNew(text string) *TextView {

	var textLines = strings.Split(text, "\n")

	var textView = TextView{text: textLines}

	textView.Init()

	for r := range textLines {
		textView.maxWidth = uint16(util.Max(int(textView.maxWidth), len(textLines[r])))
	}

	return &textView
}

func (t *TextView) Draw() {
	var y uint16 = 0

	var firstLine = int(t.ypos)
	var lastLine = int(t.ypos + t.rect.h - 1)

	var length = len(t.text)

	if lastLine >= length {
		lastLine = length - 1
		firstLine = util.Max(0, lastLine-int(t.rect.h))
	}

	for v := firstLine; v <= lastLine; v++ {

		GotoXY(t.rect.x, t.rect.y+y)

		var line = t.text[v]
		if t.xpos < uint16(len(line)) {
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
	for ; y < t.rect.h; y++ {
		GotoXY(t.rect.x, t.rect.y+y)
		WriteFill("", t.rect.w)
	}
}

func (t *TextView) ScrollBack() {
	if t.ypos > 0 {
		t.ypos--
	}
	t.RequestRedraw()
}

func (t *TextView) ScrollFwd() {
	if t.ypos+t.rect.h < uint16(len(t.text))-1 {
		t.ypos++
	}
	t.RequestRedraw()
}

func (t *TextView) ScrollLeft() {
	if t.xpos > 0 {
		t.xpos--
	}
	t.RequestRedraw()
}

func (t *TextView) ScrollRight() {
	if t.xpos+t.rect.w < t.maxWidth {
		t.xpos++
	}
	t.RequestRedraw()
}

func (t *TextView) ScrollPageFwd() {

	var ypos int = int(t.ypos) + int(t.rect.h)

	var length = len(t.text)

	if ypos+int(t.rect.h) > length-1 {
		ypos = length - int(t.rect.h) - 1
	}

	t.ypos = uint16(ypos)
	t.RequestRedraw()
}

func (t *TextView) ScrollPageBack() {
	var ypos int = util.Max(0, int(t.ypos)-int(t.rect.h))
	t.ypos = uint16(ypos)
	t.RequestRedraw()
}

func (t *TextView) HandleInput(input input.KeyInput) {

	switch input.GetKey() {
	case keyboard.KeyArrowDown:
		t.ScrollFwd()
	case keyboard.KeyPgup:
		t.ScrollPageBack()
	case keyboard.KeyArrowUp:
		t.ScrollBack()
	case keyboard.KeyPgdn:
		t.ScrollPageFwd()
	case keyboard.KeyArrowLeft:
		t.ScrollLeft()
	case keyboard.KeyArrowRight:
		t.ScrollRight()
	default:
		t.ViewImpl.HandleInput(input)
	}

}
