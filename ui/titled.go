package ui

import (
	"github.com/clidockermgr/input"
)

/**
	A container with title
**/
type TitledContainer struct {
	ViewImpl
	title  string
	child  View
	border bool
}

func TitledContainerNew(title string, child View, border bool) *TitledContainer {
	var container = TitledContainer{title: title, child: child, border: border}
	container.Init()
	return &container
}

func (t *TitledContainer) SetRect(rect Rect) {
	t.ViewImpl.SetRect(rect)
	var padding uint8 = 0
	if t.border {
		padding = 1
	}
	t.child.SetRect(Rect{x: rect.x + padding, y: rect.y + 1, w: rect.w - (padding * 2), h: rect.h - 1 - padding})
}

func (t *TitledContainer) Draw() {
	GotoXY(t.rect.x, t.rect.y)
	Background(3)
	Foreground(0)
	WriteFill(t.title, t.rect.w)
	if t.border {
		GotoXY(t.rect.x, t.rect.y+t.rect.h-1)
		WriteFill("", t.rect.w)
		WriteV(" ", t.rect.x, t.rect.y+1, t.rect.h-1)
		WriteV(" ", t.rect.x+t.rect.w-1, t.rect.y+1, t.rect.h-1)
	}
	Reset()
	t.child.Draw()
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
