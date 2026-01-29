package keyboards

import (
	"rr/kotatsutgbot/config"

	"github.com/go-telegram/bot/models"
)

type Keyboard struct {
	buttons         [][]models.KeyboardButton
	resizeKeyboard  bool
	oneTimeKeyboard bool
}

func (k *Keyboard) OneTime() *Keyboard {
	k.oneTimeKeyboard = true
	return k
}

func (k *Keyboard) MultiTime() *Keyboard {
	k.oneTimeKeyboard = false
	return k
}

func (k *Keyboard) Resize() *Keyboard {
	k.resizeKeyboard = true
	return k
}

func (k *Keyboard) DontResize() *Keyboard {
	k.resizeKeyboard = false
	return k
}

func (k *Keyboard) Row() *Keyboard {
	k.buttons = append(k.buttons, []models.KeyboardButton{})
	return k
}

func (k *Keyboard) Text(text string) *Keyboard {
	if len(k.buttons) == 0 {
		k.buttons = append(k.buttons, []models.KeyboardButton{})
	}
	k.buttons[len(k.buttons)-1] = append(k.buttons[len(k.buttons)-1], models.KeyboardButton{
		Text: text,
	})
	return k
}

func (k *Keyboard) TextC(text string, cond bool) *Keyboard {
	if !cond {
		return k
	}
	return k.Text(text)
}

func (k *Keyboard) TextIf(cond bool, text1, text2 string) *Keyboard {
	if cond {
		return k.Text(text1)
	} else {
		return k.Text(text2)
	}
}

func (k *Keyboard) TextT(text string) *Keyboard {
	return k.Text(config.T(text))
}

func (k *Keyboard) TextTC(text string, cond bool) *Keyboard {
	if !cond {
		return k
	}
	return k.TextT(text)
}

func (k *Keyboard) TextTIf(cond bool, text1, text2 string) *Keyboard {
	if cond {
		return k.TextT(text1)
	} else {
		return k.TextT(text2)
	}
}

func (k *Keyboard) RequestContact(text string) *Keyboard {
	if len(k.buttons) == 0 {
		k.buttons = append(k.buttons, []models.KeyboardButton{})
	}
	k.buttons[len(k.buttons)-1] = append(k.buttons[len(k.buttons)-1], models.KeyboardButton{
		Text:           text,
		RequestContact: true,
	})
	return k
}

func (k *Keyboard) RequestContactT(text string) *Keyboard {
	return k.RequestContact(config.T(text))
}

func (k *Keyboard) Build() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard:        k.buttons,
		ResizeKeyboard:  k.resizeKeyboard,
		OneTimeKeyboard: k.oneTimeKeyboard,
	}
}

func Default() *Keyboard {
	return &Keyboard{
		resizeKeyboard:  true,
		oneTimeKeyboard: false,
	}
}

type InlineKeyboard struct {
	buttons [][]models.InlineKeyboardButton
}

func (k *InlineKeyboard) Row() *InlineKeyboard {
	k.buttons = append(k.buttons, []models.InlineKeyboardButton{})
	return k
}

func (k *InlineKeyboard) Data(text, callback_data string) *InlineKeyboard {
	if len(k.buttons) == 0 {
		k.buttons = append(k.buttons, []models.InlineKeyboardButton{})
	}
	k.buttons[len(k.buttons)-1] = append(k.buttons[len(k.buttons)-1], models.InlineKeyboardButton{
		Text:         text,
		CallbackData: callback_data,
	})
	return k
}

func (k *InlineKeyboard) DataC(text, data string, cond bool) *InlineKeyboard {
	if !cond {
		return k
	}
	return k.Data(text, data)
}

func (k *InlineKeyboard) DataIf(cond bool, text1, data1, text2, data2 string) *InlineKeyboard {
	if cond {
		return k.Data(text1, data1)
	} else {
		return k.Data(text2, data2)
	}
}

func (k *InlineKeyboard) DataT(text, data string) *InlineKeyboard {
	return k.Data(config.T(text), data)
}

func (k *InlineKeyboard) DataTC(text, data string, cond bool) *InlineKeyboard {
	if !cond {
		return k
	}
	return k.DataT(text, data)
}

func (k *InlineKeyboard) DataTIf(cond bool, text1, data1, text2, data2 string) *InlineKeyboard {
	if cond {
		return k.DataT(text1, data1)
	} else {
		return k.DataT(text2, data2)
	}
}

func (k *InlineKeyboard) Build() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: k.buttons,
	}
}

func DefaultInline() *InlineKeyboard {
	return &InlineKeyboard{}
}
