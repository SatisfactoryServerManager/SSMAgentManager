package customwidgets

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type numericalEntry struct {
	widget.Entry
}

func NewNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	entry.SetValue(0)
	return entry
}

func (e *numericalEntry) SetValue(value int) {
	e.Entry.Text = strconv.Itoa(value)
}

func (e *numericalEntry) ParseTextValue(value string) {
	val, _ := strconv.Atoi(value)
	e.SetValue(val)
}

func (e *numericalEntry) TypedRune(r rune) {
	if r >= '0' && r <= '9' {
		e.Entry.TypedRune(r)
	}
}

func (e *numericalEntry) FocusGained() {
	e.Entry.FocusGained()
	if e.Entry.Text == "" {
		e.SetValue(0)
	} else {
		e.ParseTextValue(e.Entry.Text)
	}
}

func (e *numericalEntry) FocusLost() {
	if e.Entry.Text == "" {
		e.SetValue(0)
	} else {
		e.ParseTextValue(e.Entry.Text)
	}
	e.Entry.FocusLost()
}

func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *numericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}
