package input

import "github.com/eiannone/keyboard"

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

func GetKeyInput() (KeyInput, error) {
	input, key, err := keyboard.GetSingleKey()

	if err != nil {
		return KeyInput{}, err
	}

	return KeyInput{key: key, char: input}, nil
}
