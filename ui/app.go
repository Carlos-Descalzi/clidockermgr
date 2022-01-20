package ui

import (
	"container/list"
	"log"
	"time"

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
	log.Print("Checking input \n")
	input, key, err := keyboard.GetSingleKey()

	if err != nil {
		panic(err)
	}
	switch key {
	case keyboard.KeyTab:
		a.CycleCurrent()
	case keyboard.KeyEsc:
		a.running = false
	default:
		log.Printf("key typed: %c %d\n", input, key)
		if a.currentPopup != nil {
			a.currentPopup.HandleInput(key)
		} else {
			var currentView = a.CurrentView()
			if currentView != nil {
				currentView.HandleInput(key)
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
