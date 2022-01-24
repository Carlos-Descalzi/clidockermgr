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
	SetProperty(property int, value interface{})
	Update()
}

/**
	A default empty list model
**/
type EmptyModel struct{}

func (m *EmptyModel) ItemCount() int {
	return 0
}

func (m *EmptyModel) Item(index int) ListItem {
	return nil
}

func (m *EmptyModel) AddListener(listener ListModelListener) {}

func (m *EmptyModel) SetProperty(property int, value interface{}) {}

func (m *EmptyModel) Update() {}

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

func (m *BaseListModel) SetProperty(property int, value interface{}) {}

func (m *BaseListModel) Update() {}

/**
	A list component
**/
type List struct {
	ViewImpl
	Model         ListModel
	startIndex    int
	selectedIndex int
}

func ListNew() *List {
	var list = List{Model: &EmptyModel{}}
	list.Init()
	return &list
}

func (l *List) SetModel(model ListModel) {
	l.Model = model
	if l.Model != nil {
		l.Model.AddListener(l.Changed)
	}
}

func (l *List) Update() {
	l.Model.Update()
	if l.selectedIndex > l.Model.ItemCount() {
		l.selectedIndex = l.Model.ItemCount() - 1
	}
}

func (l *List) SelectedItem() ListItem {
	return l.Model.Item(l.selectedIndex)
}

func (l *List) Draw() {
	GotoXY(l.rect.x, l.rect.y)

	var y uint8 = 0

	for i := l.startIndex; i < l.Model.ItemCount() && y <= l.rect.h; i++ {
		GotoXY(l.rect.x, l.rect.y+y)
		var text = l.Model.Item(i)
		if l.focused && l.selectedIndex == i {
			UnderlineOn()
		}
		WriteFill(text.String(), l.rect.w)
		Reset()
		y++
	}
	if y < l.rect.h {
		for ; y <= l.rect.h; y++ {
			GotoXY(l.rect.x, l.rect.y+y)
			WriteFill("", l.rect.w)
		}
	}
}

func (l *List) ScrollFwd() {
	if l.selectedIndex < l.Model.ItemCount()-1 {
		l.selectedIndex++

		if l.selectedIndex-l.startIndex > int(l.rect.h) {
			l.startIndex++
		}
	}
}

func (l *List) ScrollBack() {
	if l.selectedIndex > 0 {
		l.selectedIndex--
	}
	if l.startIndex > l.selectedIndex {
		l.startIndex = l.selectedIndex
	}
}

func (l *List) HandleInput(input input.KeyInput) {

	switch input.GetKey() {
	case keyboard.KeyArrowDown:
		l.ScrollFwd()
		l.RequestRedraw()
	case keyboard.KeyArrowUp:
		l.ScrollBack()
		l.RequestRedraw()
	default:
		l.ViewImpl.HandleInput(input)
	}
}

func (l *List) Changed() {
	l.RequestRedraw()
}
