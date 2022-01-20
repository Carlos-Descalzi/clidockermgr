package ui

import (
	"github.com/eiannone/keyboard"
)

type Rect struct {
	x uint8
	y uint8
	w uint8
	h uint8
}

type KeyHandler func(keyboard.Key)

func RectNew(x uint8, y uint8, w uint8, h uint8) Rect {
	return Rect{x, y, w, h}
}

/**
	Interface to be implemented by view components
**/
type View interface {
	SetRect(rect Rect)
	SetVisible(visible bool)
	SetFocusable(focusable bool)
	HandleInput(key keyboard.Key)
	SetFocused(focused bool)
	AddKeyHandler(key keyboard.Key, handler KeyHandler)
	IsFocusable() bool
	Draw()
}

/**
	Base struct for views
**/
type ViewImpl struct {
	rect      Rect
	visible   bool
	focusable bool
	focused   bool
	handlers  map[keyboard.Key]KeyHandler
}

func (v *ViewImpl) Init() {
	v.handlers = make(map[keyboard.Key]KeyHandler)
}

func (v *ViewImpl) SetRect(rect Rect) {
	v.rect = rect
}

func (v *ViewImpl) SetFocusable(focusable bool) {
	v.focusable = focusable
}

func (v *ViewImpl) SetVisible(visible bool) {
	v.visible = visible
}

func (v *ViewImpl) HandleInput(key keyboard.Key) {
	var handler = v.handlers[key]

	if handler != nil {
		handler(key)
	}
}

func (v *ViewImpl) SetFocused(focused bool) {
	v.focused = focused
}

func (v ViewImpl) IsFocusable() bool {
	return v.focusable
}

func (v *ViewImpl) AddKeyHandler(key keyboard.Key, handler KeyHandler) {
	v.handlers[key] = handler
}
