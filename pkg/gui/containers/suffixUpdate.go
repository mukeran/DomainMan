package containers

import (
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	"DomainMan/pkg/whois"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SuffixUpdate(inst *instance.Instance, thisWindow fyne.Window, list *widget.List) fyne.CanvasObject {
	override := widget.NewCheck("Override Existed?", nil)
	return container.NewVBox(
		widget.NewLabel("Please select a update mode"),
		override,
		widget.NewButton("Using GitHub rfc1036/whois", func() {
			if thisWindow != nil {
				thisWindow.Close()
			}
			content := container.NewVBox(widget.NewLabel("Updating..."))
			window := utils.ShowSimpleWindow(inst, "Updating Suffixes Using GitHub rfc1036/whois", content)
			go func() {
				err := whois.UpdateServersWithRfc1036Whois(override.Checked)
				if err != nil {
					content.Add(widget.NewLabel(fmt.Sprintf("Error when updating suffixes: %v", err)))
					content.Add(widget.NewButton("Close", func() {
						window.Close()
					}))
					return
				}
				if list != nil {
					err = loadSuffixes()
					if err != nil {
						utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(fmt.Sprintf("Error when loading suffixes: %v", err)))
					}
					list.Refresh()
				}
				window.Close()
			}()
		}),
		widget.NewButton("Using Preset", func() {
			if thisWindow != nil {
				thisWindow.Close()
			}
			content := container.NewVBox(widget.NewLabel("Updating..."))
			window := utils.ShowSimpleWindow(inst, "Updating Suffixes Using Preset", content)
			go func() {
				err := whois.UpdateServersWithPreset(override.Checked)
				if err != nil {
					content.Add(widget.NewLabel(fmt.Sprintf("Error when updating suffixes: %v", err)))
					content.Add(widget.NewButton("Close", func() {
						window.Close()
					}))
					return
				}
				if list != nil {
					err = loadSuffixes()
					if err != nil {
						utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(fmt.Sprintf("Error when loading suffixes: %v", err)))
					}
					list.Refresh()
				}
				window.Close()
			}()
		}),
	)
}
