package gt_label

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"image/color"
)

type GTLabel struct {
	widget.BaseWidget
	Id        string
	Text      string
	TextColor color.Color
	Size      fyne.Size
}

// NewGTLabel 请使用这种方式初始化，直接用结构体初始化不会产生 Id 等扩展字段
func NewGTLabel(text string, size fyne.Size, textColor color.Color) *GTLabel {
	alabel := &GTLabel{
		Text:      text,
		TextColor: textColor,
	}

	alabel.Id = uuid.New().String()
	alabel.Resize(size)
	alabel.ExtendBaseWidget(alabel)
	return alabel
}

func (l *GTLabel) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(l.Text, l.TextColor)
	text.Alignment = fyne.TextAlignCenter
	return &gtLabelRenderer{label: l, text: text, objects: []fyne.CanvasObject{text}}
}

type gtLabelRenderer struct {
	label   *GTLabel
	text    *canvas.Text
	objects []fyne.CanvasObject
}

func (r *gtLabelRenderer) Layout(size fyne.Size) {
	r.text.Resize(size)
}

func (r *gtLabelRenderer) MinSize() fyne.Size {
	return r.label.Size // 返回自定义大小
}

func (r *gtLabelRenderer) Refresh() {
	r.text.Text = r.label.Text
	r.text.Color = r.label.TextColor
	r.text.Refresh()
}

func (r *gtLabelRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *gtLabelRenderer) Destroy() {}
