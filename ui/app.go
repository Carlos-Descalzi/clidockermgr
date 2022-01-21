package ui

import (
	"container/list"
	"time"

	"github.com/clidockermgr/input"
	"github.com/eiannone/keyboard"
)

type Application struct {
	children       *list.List
	currentElement *list.Element
	running        bool
	currentPopup   View
}

func ApplicationNew() *Application {
	return &Application{children: list.New(), running: true, currentElement: nil, currentPopup: nil}
}

func (a *Application) Add(view View) {
	var empty = a.children.Len() == 0
	a.children.PushBack(view)
	if empty {
		a.currentElement = a.children.Front()
	}
	view.AddRedrawListener(a.RedrawRequested)
}

func (a *Application) ShowPopup(view View) {
	a.currentPopup = view
}

func (a *Application) ClosePopup() {
	a.currentPopup = nil
}

func (a *Application) CycleCurrent() {
	if a.currentElement != nil {
		a.currentElement.Value.(View).SetFocused(false)
	}
	a.currentElement = a.currentElement.Next()
	if a.currentElement == nil {
		a.currentElement = a.children.Front()
	}
	a.currentElement.Value.(View).SetFocused(true)
}

func (a *Application) CurrentView() View {
	if a.currentElement == nil {
		return nil
	}
	return (a.currentElement.Value.(View))
}

func (a *Application) CheckInput() {
	input, err := input.GetKeyInput()

	if err != nil {
		panic(err)
	}
	key := input.GetKey()
	switch key {
	case keyboard.KeyTab:
		a.CycleCurrent()
	case keyboard.KeyEsc:
		if a.currentPopup != nil {
			a.ClosePopup()
		} else {
			a.running = false
		}
	default:
		if a.currentPopup != nil {
			a.currentPopup.HandleInput(input)
		} else {
			var currentView = a.CurrentView()
			if currentView != nil {
				currentView.HandleInput(input)
			}
		}
	}

}

func (a *Application) DrawAll() {
	if a.currentPopup != nil {
		a.currentPopup.Draw()
	} else {
		for v := a.children.Front(); v != nil; v = v.Next() {
			var view View = v.Value.(View)
			if view != nil {
				view.Draw()
			}
		}
	}
}

func (a *Application) RedrawRequested(view interface{}) {
	// TODO: only redraw component
	a.DrawAll()
}

func (a *Application) Loop() {
	keyboard.Open()
	ClearScreen()
	CursorOff()
	a.DrawAll()
	for a.running {
		a.CheckInput()
		a.DrawAll()
		time.Sleep(time.Duration(10))
	}
	CursorOn()
	ClearScreen()
}
