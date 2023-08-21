package gt_button

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"image/color"
)

type GTButton struct {
	widget.BaseWidget
	Id              string
	Text            string
	BackgroundColor color.Color
	TextColor       color.Color
	OnTapped        func(btn *GTButton)
}

// NewGTButton 请使用这种方式初始化，直接用结构体初始化不会产生 Id 等扩展字段
func NewGTButton(text string, size fyne.Size, backgroundColor, textColor color.Color, onTapped func(btn *GTButton)) *GTButton {
	button := &GTButton{
		Text:            text,
		BackgroundColor: backgroundColor,
		TextColor:       textColor,
		OnTapped:        onTapped,
	}

	button.Id = uuid.New().String()
	button.ExtendBaseWidget(button)

	button.Resize(size)
	return button
}

func (b *GTButton) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(b.Text, b.TextColor)
	text.Alignment = fyne.TextAlignCenter
	bg := canvas.NewRectangle(b.BackgroundColor)
	return &gtButtonRenderer{button: b, text: text, bg: bg, objects: []fyne.CanvasObject{bg, text}}
}

type gtButtonRenderer struct {
	button  *GTButton
	text    *canvas.Text
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *gtButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.text.Resize(size)
}

func (r *gtButtonRenderer) MinSize() fyne.Size {
	return r.button.Size() // 返回自定义大小
}

func (r *gtButtonRenderer) Refresh() {
	r.text.Text = r.button.Text
	r.text.Color = r.button.TextColor
	r.bg.FillColor = r.button.BackgroundColor
	r.text.Refresh()
	r.bg.Refresh()
}

func (r *gtButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *gtButtonRenderer) Destroy() {}

func (b *GTButton) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped(b)
	}
}
