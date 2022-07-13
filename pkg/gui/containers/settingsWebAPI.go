package containers

import (
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsWebAPIWindow(inst *instance.Instance) fyne.Window {
	window := inst.App.NewWindow("Database Settings")
	window.SetContent(SettingsWebAPI(inst, window))
	window.Resize(fyne.NewSize(400, 0))
	window.CenterOnScreen()
	return window
}

func SettingsWebAPI(inst *instance.Instance, thisWindow fyne.Window) fyne.CanvasObject {
	listen := widget.NewEntry()
	listen.SetText(inst.Settings.WebAPI.Listen)
	listen.SetPlaceHolder("0.0.0.0:8899")
	var form *widget.Form
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Listen Address", Widget: listen},
		},
		OnSubmit: func() {
			inst.Settings.WebAPI.Listen = listen.Text
			err := inst.Settings.Save()
			if err != nil {
				utils.ShowMessage(inst, "Settings File Error", fmt.Sprintf("Error when saving settings: %v", err))
			}
			if thisWindow != nil {
				thisWindow.Close()
			}
		},
	}
	return form
}
