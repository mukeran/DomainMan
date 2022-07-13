package containers

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	"DomainMan/pkg/models"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SuffixDeleteAll(inst *instance.Instance, thisWindow fyne.Window, list *widget.List) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabel("Are you sure you want to delete all suffixes?"),
		widget.NewButton("Yes", func() {
			if thisWindow != nil {
				thisWindow.Close()
			}
			go func() {
				if v := database.DB.Delete(&models.Suffix{}, "1=1"); v.Error != nil {
					utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(fmt.Sprintf("Error when deleting suffixes: %v", v.Error)))
					return
				}
				if list != nil {
					err := loadSuffixes()
					if err != nil {
						utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(fmt.Sprintf("Error when loading suffixes: %v", err)))
					}
					list.Refresh()
				}
			}()
		}),
		widget.NewButton("No", func() {
			if thisWindow != nil {
				thisWindow.Close()
			}
		}),
	)
}
