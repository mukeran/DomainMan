package containers

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsDatabaseWindow(inst *instance.Instance, showMainWindow bool) fyne.Window {
	window := inst.App.NewWindow("Database Settings")
	window.SetContent(SettingsDatabase(inst, window, showMainWindow))
	window.Resize(fyne.NewSize(400, 0))
	window.CenterOnScreen()
	return window
}

func SettingsDatabase(inst *instance.Instance, thisWindow fyne.Window, showMainWindow bool) fyne.CanvasObject {
	databaseDialect := widget.NewSelect([]string{"mysql", "sqlite"}, nil)
	if inst.Settings.Database.Dialect != "" {
		databaseDialect.SetSelected(inst.Settings.Database.Dialect)
	} else {
		databaseDialect.SetSelected("sqlite")
	}
	databaseParameter := widget.NewEntry()
	if inst.Settings.Database.Parameter != "" {
		databaseParameter.SetText(inst.Settings.Database.Parameter)
	} else {
		databaseParameter.SetText("db.sqlite3")
	}
	databaseParameter.SetPlaceHolder("e.g.: db.sqlite3, mysql://root:root@127.0.0.1:3306/DomainMan")
	var form *widget.Form
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Database Dialect", Widget: databaseDialect},
			{Text: "Database Parameter", Widget: databaseParameter},
		},
		OnSubmit: func() {
			form.SubmitText = "Checking connection..."
			form.Refresh()
			form.Disable()
			err := database.Connect(databaseDialect.Selected, databaseParameter.Text, true)
			if err != nil {
				form.SubmitText = "Submit"
				form.Refresh()
				form.Enable()
				utils.ShowMessage(inst, "Database Error", fmt.Sprintf("Error when connecting to database: %v", err))
				return
			}
			inst.Settings.Database.Dialect = databaseDialect.Selected
			inst.Settings.Database.Parameter = databaseParameter.Text
			err = inst.Settings.Save()
			if err != nil {
				utils.ShowMessage(inst, "Settings File Error", fmt.Sprintf("Error when saving settings: %v", err))
			}
			if showMainWindow {
				if inst.MainWindow == nil {
					InitMainWindow(inst)
				}
				inst.MainWindow.Show()
			}
			if thisWindow != nil {
				thisWindow.Close()
			}
		},
	}
	return form
}
