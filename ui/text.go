package ui

import (
	"log"
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
	log.Printf("Text lines: %d\n", len(textView.text))
	textView.Init()

	for r := range textView.text {
		textView.maxWidth = uint8(util.Max(int(textView.maxWidth), len(textView.text[r])))
	}

	return &textView
}

func (t *TextView) Draw() {
	var y uint8 = 0

	var firstLine = int(t.ypos)
	var lastLine = int(t.ypos + t.rect.h - 1)

	var length = len(t.text)

	if lastLine >= length {
		lastLine = length - 1
		firstLine = util.Max(0, lastLine-int(t.rect.h))
	}
	log.Printf("%d %d\n", firstLine, lastLine)

	for v := firstLine; v <= lastLine; v++ {

		GotoXY(t.rect.x, t.rect.y+y)

		var line = t.text[v]
		if t.xpos < uint8(len(line)) {
			line = line[t.xpos:]
		} else {
			line = ""
		}
		log.Printf("Writing %s\n", line)

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
