package containers

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	"DomainMan/pkg/models"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	suffixes []*models.Suffix
)

func loadSuffixes() error {
	suffixes = []*models.Suffix{}
	if v := database.DB.Find(&suffixes); v.Error != nil {
		return v.Error
	}
	return nil
}

func suffixListTable(inst *instance.Instance) *widget.List {
	err := loadSuffixes()
	if err != nil {
		utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(fmt.Sprintf("Error when loading suffixes: %v", err)))
	}
	db := database.DB
	var table *widget.List
	table = widget.NewList(
		func() int {
			return len(suffixes)
		},
		func() fyne.CanvasObject {
			nameLabel := widget.NewLabelWithStyle(".com", fyne.TextAlignLeading, fyne.TextStyle{
				Bold: true,
			})
			nameLabel.Wrapping = 20
			return container.NewHBox(
				nameLabel,
				layout.NewSpacer(),
				widget.NewLabel("whois.verisign-grs.com:43"),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(suffixes[i].Name)
			c.Objects[2].(*widget.Label).SetText(suffixes[i].WhoisServer)
			c.Objects[4].(*widget.Button).OnTapped = func(self *widget.Button) func() {
				return func() {
					self.Disable()
					if v := db.Delete(suffixes[i]); v.Error != nil {
						utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(v.Error.Error()))
					} else {
						suffixes = append(suffixes[:i], suffixes[i+1:]...)
						table.Refresh()
					}
					self.Enable()
				}
			}(c.Objects[4].(*widget.Button))
		})
	return table
}

func suffixListToolbar(inst *instance.Instance, table *widget.List) *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			utils.ShowSizedWindowComputed(inst, "Add Suffix", func(window fyne.Window) fyne.CanvasObject {
				return SuffixAdd(inst, window, table)
			}, 400, 0)
		}),
		widget.NewToolbarAction(theme.DownloadIcon(), func() {
			utils.ShowSizedWindowComputed(inst, "Update Suffixes", func(window fyne.Window) fyne.CanvasObject {
				return SuffixUpdate(inst, window, table)
			}, 400, 0)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			utils.ShowSizedWindowComputed(inst, "Delete All Suffixes", func(window fyne.Window) fyne.CanvasObject {
				return SuffixDeleteAll(inst, window, table)
			}, 400, 0)
		}),
	)
}

func SuffixList(inst *instance.Instance) fyne.CanvasObject {
	table := suffixListTable(inst)
	toolbar := suffixListToolbar(inst, table)
	return container.NewBorder(toolbar, nil, nil, nil, table)
}
