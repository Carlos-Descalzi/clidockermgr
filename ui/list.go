package ui

import (
	"container/list"
	"fmt"

	"github.com/clidockermgr/input"
	"github.com/eiannone/keyboard"
)

/**
	Interface for list items
**/
type ListItem interface {
	fmt.Stringer
	Value() interface{}
}

type ListModelListener func()

/**
	Model interface for list component
**/
type ListModel interface {
	ItemCount() int
	Item(index int) ListItem
	AddListener(listener ListModelListener)
}

/**
	A default empty list model
**/
type EmptyModel struct {
}

func (m EmptyModel) ItemCount() int {
	return 0
}

func (m EmptyModel) Item(index int) ListItem {
	return nil
}

func (m EmptyModel) AddListener(listener ListModelListener) {

}

/**
	A base list model struct which implements
	listener handling
**/
type BaseListModel struct {
	listeners *list.List
}

func (m *BaseListModel) Init() {
	m.listeners = list.New()
}

func (m *BaseListModel) AddListener(listener ListModelListener) {
	m.listeners.PushBack(listener)
}

func (m BaseListModel) NotifyChanged() {
	for v := m.listeners.Front(); v != nil; v = v.Next() {
		v.Value.(ListModelListener)()
	}
}

/**
	A list component
**/
type List struct {
	ViewImpl
	model         ListModel
	startIndex    int
	selectedIndex int
}

func ListNew() *List {
	var list = List{model: &EmptyModel{}}
	list.Init()
	return &list
}

func (l *List) SetModel(model ListModel) {
	l.model = model
	if l.model != nil {
		l.model.AddListener(l.Changed)
	}
}

func (l *List) SelectedItem() ListItem {
	return l.model.Item(l.selectedIndex)
}

func (l *List) Draw() {
	GotoXY(l.rect.x, l.rect.y)

	var y uint8 = 0

	for i := l.startIndex; i < l.model.ItemCount(); i++ {
		GotoXY(l.rect.x, l.rect.y+y)
		var text = l.model.Item(i)
		if l.focused && l.selectedIndex == i {
			UnderlineOn()
		}
		WriteFill(text.String(), l.rect.w)
		Reset()
		y++
		if y >= l.rect.h {
			break
		}
	}
	if y < l.rect.h {
		for ; y < l.rect.h; y++ {
			GotoXY(l.rect.x, l.rect.y+y)
			WriteFill("", l.rect.w)
		}
	}
}

func (l *List) ScrollBack() {
	if l.selectedIndex < l.model.ItemCount()-1 {
		l.selectedIndex++
	}
}

func (l *List) ScrollFwd() {
	if l.selectedIndex > 0 {
		l.selectedIndex--
	}
	if l.startIndex > 0 {
		l.startIndex--
	}
}

func (l *List) HandleInput(input input.KeyInput) {

	switch input.GetKey() {
	case keyboard.KeyArrowDown:
		l.ScrollBack()
	case keyboard.KeyArrowUp:
		l.ScrollFwd()
	default:
		l.ViewImpl.HandleInput(input)
	}
}

func (l *List) Changed() {
	l.RequestRedraw()
}
