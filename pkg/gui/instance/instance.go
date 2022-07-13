package instance

import (
	"fyne.io/fyne/v2"
)

type Instance struct {
	App        fyne.App
	MainWindow fyne.Window
	Settings   *Settings
}
