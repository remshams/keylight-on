package kl_textinput

import (
	"fmt"
	"keylight-charm/styles"

	"github.com/charmbracelet/bubbles/textinput"
)

func CreateTextInputModel() textinput.Model {
	model := textinput.New()
	model.TextStyle = styles.TextAccentColor
	return model
}

func CreateTextInputView(model textinput.Model, label string, unit string) string {
	return fmt.Sprintf("%s %s%s", label, model.View(), styles.TextAccentColor.Render(unit))
}
