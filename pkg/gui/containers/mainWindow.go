package containers

import (
	"DomainMan/pkg/gui/instance"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func InitMainWindow(inst *instance.Instance) {
	inst.MainWindow = inst.App.NewWindow("DomainMan GUI")

	tabs := container.NewAppTabs(
		container.NewTabItem("Domains", DomainList(inst)),
		container.NewTabItem("Whois", WhoisQuery(inst)),
		container.NewTabItem("Suffixes", SuffixList(inst)),
		container.NewTabItem("Settings", SettingsRoot(inst)),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	inst.MainWindow.Resize(fyne.NewSize(800, 400))
	inst.MainWindow.CenterOnScreen()
	inst.MainWindow.SetContent(tabs)
	inst.MainWindow.SetMaster()
}
