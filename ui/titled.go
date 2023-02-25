package ui

import (
	"fmt"
	"strings"

	"github.com/clidockermgr/input"
)

type BorderStyle func(string, Rect)

func HeaderBorder(title string, rect Rect) {
	GotoXY(rect.x, rect.y)
	Background(7)
	Foreground(0)
	WriteFill(title, rect.w)
	Reset()
}

func FullBorder(title string, rect Rect) {
	GotoXY(rect.x, rect.y)
	Background(7)
	Foreground(0)
	WriteFill(title, rect.w)
	GotoXY(rect.x, rect.y+rect.h-1)
	WriteFill("", rect.w)
	WriteV(" ", rect.x, rect.y+1, rect.h-1)
	WriteV(" ", rect.x+rect.w-1, rect.y+1, rect.h-1)
	Reset()
}

func LineBorder(title string, rect Rect) {
	GotoXY(rect.x+1, rect.y)
	fmt.Print(title)
	GotoXY(rect.x, rect.y)
	fmt.Print(LineBorderTopLeft)
	WriteV(LineBorderHorizontal, rect.x, rect.y+1, rect.h-2)
	WriteV(LineBorderHorizontal, rect.x+rect.w-1, rect.y+1, rect.h-2)
	GotoXY(rect.x, rect.y+rect.h-1)
	fmt.Print(LineBorderBottomLeft)
	GotoXY(rect.x+rect.w-1, rect.y+rect.h-1)
	fmt.Print(LineBorderBottomRight)
	GotoXY(rect.x+rect.w-1, rect.y)
	fmt.Print(LineBorderTopRight)
	GotoXY(rect.x+1, rect.y+rect.h-1)
	fmt.Print(strings.Repeat(LineBorderVertical, int(rect.w-2)))
	GotoXY(rect.x+1+uint16(len(title)), rect.y)
	fmt.Print(strings.Repeat(LineBorderVertical, int(rect.w-2)-len(title)))
}

/**
	A container with title
**/
type TitledContainer struct {
	ViewImpl
	title  string
	child  View
	Border BorderStyle
}

func TitledContainerNew(title string, child View, border bool) *TitledContainer {
	var container = TitledContainer{title: title, child: child, Border: HeaderBorder}
	container.Init()
	return &container
}

func (t *TitledContainer) SetRect(rect Rect) {
	t.ViewImpl.SetRect(rect)

	var padding uint16 = 0

	if t.Border != nil {
		padding = 1
	}

	t.child.SetRect(
		Rect{
			x: rect.x + padding,
			y: rect.y + 1,
			w: rect.w - (padding * 2),
			h: rect.h - 1 - padding})
}

func (t *TitledContainer) Draw() {
	if t.Border != nil {
		t.Border(t.title, t.rect)
	}
	t.child.Draw()
}

func (t *TitledContainer) CheckRedrawFlag() bool {
	return t.ViewImpl.CheckRedrawFlag() || t.child.CheckRedrawFlag()
}

func (t *TitledContainer) SetFocusable(focusable bool) {
	t.child.SetFocusable(focusable)
}

func (t *TitledContainer) HandleInput(input input.KeyInput) {
	t.child.HandleInput(input)
}

func (t *TitledContainer) SetFocused(focused bool) {
	t.child.SetFocused(focused)
}

func (t TitledContainer) IsFocusable() bool {
	return t.child.IsFocusable()
}
