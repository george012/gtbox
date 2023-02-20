package gtbox_gui

import (
	"github.com/george012/gtbox/gtbox_gui/gtbox_gui_fyne"
	"github.com/george012/gtbox/gtbox_gui/gtbox_gui_wails"
)

type GTBoxGUIFlag int

const (
	GTBoxGUIFlag_Fyne GTBoxGUIFlag = iota
	GTBoxGUIFlag_Wails
	GTBoxGUIFlag_None
)

func (aFlag GTBoxGUIFlag) String() string {
	switch aFlag {
	case GTBoxGUIFlag_Fyne:
		return "Fyne"
	case GTBoxGUIFlag_Wails:
		return "Wails"
	case GTBoxGUIFlag_None:
		return "未知类型"
	default:
		return "未知类型"
	}
}

// GetDescription 获取描述
func (aFlag GTBoxGUIFlag) GetDescription() string {
	switch aFlag {
	case GTBoxGUIFlag_Fyne:
		return gtbox_gui_fyne.GTGetDescription()
	case GTBoxGUIFlag_Wails:
		return gtbox_gui_wails.GTGetDescription()
	case GTBoxGUIFlag_None:
		return ""
	default:
		return ""
	}
}
