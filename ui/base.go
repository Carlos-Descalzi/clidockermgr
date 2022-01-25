package ui

import "github.com/clidockermgr/input"

type Rect struct {
	x uint8
	y uint8
	w uint8
	h uint8
}

type Insets struct {
	Top    uint8
	Bottom uint8
	Left   uint8
	Right  uint8
}

type KeyHandler func(input.KeyInput)
type RedrawListener func(view interface{})

func RectNew(x, y, w, h uint8) Rect {
	return Rect{x, y, w, h}
}
