package ui

import "github.com/clidockermgr/input"

type Rect struct {
	x uint16
	y uint16
	w uint16
	h uint16
}

type Insets struct {
	Top    uint16
	Bottom uint16
	Left   uint16
	Right  uint16
}

type KeyHandler func(input.KeyInput)
type RedrawListener func(view interface{})

func RectNew(x, y, w, h uint16) Rect {
	return Rect{x, y, w, h}
}
