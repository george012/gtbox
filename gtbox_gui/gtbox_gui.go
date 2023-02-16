package gtbox_gui

import (
	"github.com/george012/gtbox/gtbox_gui/gtbox_gui_fyne"
	"github.com/george012/gtbox/gtbox_gui/gtbox_gui_wails"
)

type GTGOGUIFlag int

const (
	GTGOGUIFlag_Fyne GTGOGUIFlag = iota
	GTGOGUIFlag_Wails
	GTGOGUIFlag_None
)

func (aFlag GTGOGUIFlag) String() string {
	switch aFlag {
	case GTGOGUIFlag_Fyne:
		return "Fyne"
	case GTGOGUIFlag_Wails:
		return "Wails"
	case GTGOGUIFlag_None:
		return "未知类型"
	default:
		return "未知类型"
	}
}

// GetDescription 获取描述
func (aFlag GTGOGUIFlag) GetDescription() string {
	switch aFlag {
	case GTGOGUIFlag_Fyne:
		return gtbox_gui_fyne.GTGetDescription()
	case GTGOGUIFlag_Wails:
		return gtbox_gui_wails.GTGetDescription()
	case GTGOGUIFlag_None:
		return ""
	default:
		return ""
	}
}
