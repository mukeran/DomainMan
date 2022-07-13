package gui

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/gui/containers"
	"DomainMan/pkg/gui/instance"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
)

func connectDatabase(inst *instance.Instance) fyne.Window {
	if os.Getenv("DOMAINMAN_DATABASE_DIALECT") != "" {
		inst.Settings.Database.Dialect = os.Getenv("DOMAINMAN_DATABASE_DIALECT")
	}
	if os.Getenv("DOMAINMAN_DATABASE_PARAMETER") != "" {
		inst.Settings.Database.Parameter = os.Getenv("DOMAINMAN_DATABASE_PARAMETER")
	}
	if inst.Settings.Database.Dialect == "" || inst.Settings.Database.Parameter == "" {
		return containers.NewSettingsDatabaseWindow(inst, true)
	}
	err := database.Connect(inst.Settings.Database.Dialect, inst.Settings.Database.Parameter, true)
	if err != nil {
		window := inst.App.NewWindow("Database Error")
		window.SetContent(container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Error when connecting to database: %v", err)),
			container.NewVBox(
				widget.NewButton("Retry", func() {
					window.Close()
					connectDatabase(inst).Show()
				}),
				widget.NewButton("Change Settings", func() {
					window.Close()
					containers.NewSettingsDatabaseWindow(inst, true).Show()
				}),
			),
		))
		window.CenterOnScreen()
		return window
	}
	return nil
}

func Run() {
	inst := &instance.Instance{App: app.New()}
	inst.App.Settings().SetTheme(&myTheme{})

	var err error
	inst.Settings, err = instance.LoadSettings()
	if err != nil {
		inst.Settings = instance.GetDefaultSettings()
	}

	window := connectDatabase(inst)
	if window != nil {
		window.ShowAndRun()
	} else {
		containers.InitMainWindow(inst)
		inst.MainWindow.ShowAndRun()
	}
}
