package input

import (
	"log"

	"github.com/eiannone/keyboard"
)

type KeyInput struct {
	key  keyboard.Key
	char rune
}

func KeyInputChar(char rune) KeyInput {
	return KeyInput{char: char}
}

func KeyInputKey(key keyboard.Key) KeyInput {
	return KeyInput{key: key}
}

func KeyInputNew(key keyboard.Key, char rune) KeyInput {
	return KeyInput{key: key, char: char}
}

func (k KeyInput) GetKey() keyboard.Key {
	return k.key
}

func (k KeyInput) GetChar() rune {
	return k.char
}

type InputHandler struct {
	active  bool
	channel chan KeyInput
}

func InputHandlerNew() *InputHandler {
	handler := InputHandler{active: true, channel: make(chan KeyInput, 100)}
	go handler.RunCheck()
	return &handler
}

func (i *InputHandler) RunCheck() {
	keyboard.Open()
	defer keyboard.Close()
	for i.active {
		input, key, err := keyboard.GetKey()
		if err == nil {
			i.channel <- KeyInput{key: key, char: input}
		} else {
			log.Print(err)
		}
	}
}

func (i *InputHandler) GetKeyInput() (KeyInput, bool) {
	if len(i.channel) > 0 {
		return <-i.channel, true
	}
	return KeyInput{}, false
}

func (i *InputHandler) Close() {
	keyboard.Close()
	i.active = false
	close(i.channel)
}
