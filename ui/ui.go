package ui

import (
	"github.com/clidockermgr/input"
)

type Rect struct {
	x uint8
	y uint8
	w uint8
	h uint8
}

type KeyHandler func(input.KeyInput)
type RedrawListener func(view interface{})

func RectNew(x, y, w, h uint8) Rect {
	return Rect{x, y, w, h}
}

/**
	Interface to be implemented by view components
**/
type View interface {
	SetRect(rect Rect)
	SetVisible(visible bool)
	SetFocusable(focusable bool)
	HandleInput(input input.KeyInput)
	SetFocused(focused bool)
	AddKeyHandler(input input.KeyInput, handler KeyHandler)
	IsFocusable() bool
	Draw()
	CheckRedrawFlag() bool
	RequestRedraw()
}

/**
	Base struct for views
**/
type ViewImpl struct {
	rect      Rect
	visible   bool
	focusable bool
	focused   bool
	dirty     bool
	handlers  map[input.KeyInput]KeyHandler
}

func (v *ViewImpl) Init() {
	v.handlers = make(map[input.KeyInput]KeyHandler)
	v.dirty = true
}

func (v *ViewImpl) SetRect(rect Rect) {
	v.rect = rect
	v.RequestRedraw()
}

func (v *ViewImpl) SetFocusable(focusable bool) {
	v.focusable = focusable
}

func (v *ViewImpl) SetVisible(visible bool) {
	v.visible = visible
	v.RequestRedraw()
}

func (v *ViewImpl) HandleInput(input input.KeyInput) {
	var handler = v.handlers[input]

	if handler != nil {
		handler(input)
	}
}

func (v *ViewImpl) SetFocused(focused bool) {
	v.focused = focused
	v.RequestRedraw()
}

func (v ViewImpl) IsFocusable() bool {
	return v.focusable
}

func (v *ViewImpl) AddKeyHandler(key input.KeyInput, handler KeyHandler) {
	v.handlers[key] = handler
}

func (v *ViewImpl) RequestRedraw() {
	v.dirty = true
}

func (v *ViewImpl) CheckRedrawFlag() bool {
	flag := v.dirty
	v.dirty = false
	return flag
}
