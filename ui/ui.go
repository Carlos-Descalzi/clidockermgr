package ui

import (
	"container/list"

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
	HandleInput(input input.KeyInput)
	SetFocused(focused bool)
	AddKeyHandler(input input.KeyInput, handler KeyHandler)
	AddRedrawListener(listener RedrawListener)
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
	handlers  map[input.KeyInput]KeyHandler
	listeners *list.List
}

func (v *ViewImpl) Init() {
	v.handlers = make(map[input.KeyInput]KeyHandler)
	v.listeners = list.New()
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

func (v *ViewImpl) HandleInput(input input.KeyInput) {
	var handler = v.handlers[input]

	if handler != nil {
		handler(input)
	}
}

func (v *ViewImpl) SetFocused(focused bool) {
	v.focused = focused
}

func (v ViewImpl) IsFocusable() bool {
	return v.focusable
}

func (v *ViewImpl) AddKeyHandler(key input.KeyInput, handler KeyHandler) {
	v.handlers[key] = handler
}

func (v *ViewImpl) AddRedrawListener(listener RedrawListener) {
	v.listeners.PushBack(listener)
}

func (v *ViewImpl) RequestRedraw() {
	for i := v.listeners.Front(); i != nil; i = i.Next() {
		i.Value.(RedrawListener)(v)
	}
}
