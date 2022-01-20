package ui

import (
	"fmt"
	"strings"
)

/**
	A label component
**/
type Label struct {
	ViewImpl
	text string
}

func LabelNew(text string) *Label {
	var label = Label{text: text}
	label.Init()
	return &label
}

func (l *Label) Draw() {
	GotoXY(l.rect.x, l.rect.y)
	if len(l.text) > int(l.rect.w) {
		fmt.Print(l.text[0:l.rect.w])
	} else {
		fmt.Printf("%s%s", l.text, strings.Repeat(" ", (int(l.rect.w)-len(l.text))))
	}
}
