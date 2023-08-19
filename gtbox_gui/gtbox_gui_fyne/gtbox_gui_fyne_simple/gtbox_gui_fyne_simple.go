package gtbox_gui_fyne_simple

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/george012/gtbox/gtbox_gui/gtbox_gui_fyne"
	"github.com/george012/gtbox/gtbox_gui/gtbox_gui_fyne/gt_button"
	"image/color"
)

func SimpleWithGTButton() {
	myApp := app.New()
	myWindow := myApp.NewWindow("GTButton Example")

	gtButton := gt_button.NewGTButton("Click Me",
		fyne.NewSize(100, 200),
		color.NRGBA{R: 255, G: 0, B: 0, A: 255},
		color.White,
		func(btn *gt_button.GTButton) {
			fmt.Printf("click Bution With [%s]", btn.Id)
			fmt.Print("%s", gtbox_gui_fyne.GTGetDescription())
		},
	)

	myWindow.SetContent(container.NewVBox(gtButton))
	myWindow.ShowAndRun()
}
