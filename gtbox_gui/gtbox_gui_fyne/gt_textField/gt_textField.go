package gt_textField

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"image/color"
)

type GTTextField struct {
	widget.BaseWidget
	Id              string
	Placeholder     string
	BackgroundColor color.Color
	OnChanged       func(textField *GTTextField, value string)
	entry           *widget.Entry
}

func NewGTInput(placeholder string, size fyne.Size, backgroundColor color.Color, onValueChanged func(textField *GTTextField, value string)) *GTTextField {
	atextField := &GTTextField{
		Placeholder:     placeholder,
		BackgroundColor: backgroundColor,
		OnChanged:       onValueChanged,
	}

	atextField.Id = uuid.New().String()
	atextField.entry = widget.NewEntry()
	atextField.entry.SetPlaceHolder(placeholder)
	atextField.entry.OnChanged = func(text string) {
		if onValueChanged != nil {
			onValueChanged(atextField, text)
		}
	}

	atextField.ExtendBaseWidget(atextField)

	atextField.entry.Resize(size) // 设置输入框的大小

	return atextField
}

func (t *GTTextField) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(t.BackgroundColor)
	return &gtTextFieldRenderer{textField: t, bg: bg, objects: []fyne.CanvasObject{bg, t.entry}}
}

type gtTextFieldRenderer struct {
	textField *GTTextField
	bg        *canvas.Rectangle
	objects   []fyne.CanvasObject
}

func (r *gtTextFieldRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.textField.entry.Resize(size)
}

func (r *gtTextFieldRenderer) MinSize() fyne.Size {
	return r.textField.entry.MinSize() // 返回输入框的最小大小
}

func (r *gtTextFieldRenderer) Refresh() {
	r.bg.FillColor = r.textField.BackgroundColor
	r.bg.Refresh()
	r.textField.entry.Refresh()
}

func (r *gtTextFieldRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *gtTextFieldRenderer) Destroy() {}

func (t *GTTextField) SetValue(value string) {
	t.entry.SetText(value)
}

func (t *GTTextField) GetValue() string {
	return t.entry.Text
}
