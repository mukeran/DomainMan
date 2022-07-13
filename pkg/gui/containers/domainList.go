package containers

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	models2 "DomainMan/pkg/models"
	"DomainMan/pkg/whois"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

var (
	domains []*models2.Domain
	whoises []*models2.Whois
)

func loadDomains() error {
	domains = []*models2.Domain{}
	whoises = []*models2.Whois{}
	db := database.DB
	if v := db.Find(&domains); v.Error != nil {
		return v.Error
	}
	for _, domain := range domains {
		var w models2.Whois
		if v := db.Where("domain_name = ?", domain.Name).Order("created_at desc").First(&w); errors.Is(v.Error, gorm.ErrRecordNotFound) {
			whoises = append(whoises, nil)
		} else if v.Error != nil {
			return v.Error
		} else {
			whoises = append(whoises, &w)
		}
	}
	return nil
}

func domainListTable(inst *instance.Instance) *widget.List {
	err := loadDomains()
	if err != nil {
		utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(fmt.Sprintf("Error when loading domains: %v", err)))
	}
	db := database.DB
	var table *widget.List
	table = widget.NewList(
		func() int {
			return len(domains)
		},
		func() fyne.CanvasObject {
			nameLabel := widget.NewLabelWithStyle("example.com", fyne.TextAlignLeading, fyne.TextStyle{
				Bold: true,
			})
			nameLabel.Wrapping = 30
			return container.NewHBox(
				nameLabel,
				layout.NewSpacer(),
				widget.NewLabel("0000-00-00|0000-00-00"),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("", theme.VisibilityIcon(), func() {
				}),
				widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
				}),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(domains[i].Name)
			if whoises[i] != nil {
				expirationDate := whoises[i].ExpirationDate.Format("2006-01-02")
				registrationDate := whoises[i].RegistrationDate.Format("2006-01-02")
				if whoises[i].RegistrationDate.IsZero() {
					registrationDate = "-REDACTED-"
				}
				c.Objects[2].(*widget.Label).SetText(fmt.Sprintf("%s|%s", expirationDate, registrationDate))
			} else {
				c.Objects[2].(*widget.Label).SetText("not found")
			}
			if whoises[i] != nil {
				c.Objects[4].(*widget.Button).Enable()
				c.Objects[4].(*widget.Button).OnTapped = func() {
					utils.ShowSizedWindow(inst, fmt.Sprintf("Whois Info of %s", domains[i].Name), WhoisDetail(inst, whoises[i]), 600, 400)
				}
			} else {
				c.Objects[4].(*widget.Button).Disable()
			}
			c.Objects[5].(*widget.Button).OnTapped = func(self *widget.Button) func() {
				return func() {
					self.Disable()
					self.SetIcon(theme.MoreHorizontalIcon())
					go func() {
						w, err := whois.Lookup(domains[i].Name)
						if err != nil {
							utils.ShowSimpleWindow(inst, "Whois Lookup Error", widget.NewLabel(err.Error()))
						} else {
							if v := db.Create(w); v.Error != nil {
								utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(v.Error.Error()))
							}
						}
						self.Enable()
						self.SetIcon(theme.ViewRefreshIcon())
						if err == nil {
							whoises[i] = w
							c.Objects[2].(*widget.Label).SetText(fmt.Sprintf("%s|%s", whoises[i].ExpirationDate.Format("2006-01-02"), whoises[i].RegistrationDate.Format("2006-01-02")))
						}
					}()
				}
			}(c.Objects[5].(*widget.Button))
			c.Objects[6].(*widget.Button).OnTapped = func(self *widget.Button) func() {
				return func() {
					self.Disable()
					if v := db.Delete(domains[i]); v.Error != nil {
						utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(v.Error.Error()))
					} else {
						domains = append(domains[:i], domains[i+1:]...)
						whoises = append(whoises[:i], whoises[i+1:]...)
						table.Refresh()
					}
					self.Enable()
				}
			}(c.Objects[6].(*widget.Button))
		})
	return table
}

func domainListToolbar(inst *instance.Instance, table *widget.List) *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			utils.ShowSizedWindowComputed(inst, "Add Domain", func(window fyne.Window) fyne.CanvasObject {
				return DomainAdd(inst, window, table)
			}, 400, 0)
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			err := loadDomains()
			if err != nil {
				utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(fmt.Sprintf("Error when loading domains: %v", err)))
				return
			}
			content := container.NewVBox(
				widget.NewLabel(fmt.Sprintf("Ready to refresh domains' WHOIS info. [0/%d]", len(domains))),
			)
			window := utils.ShowSizedWindow(inst, fmt.Sprintf("Refreshing... [0/%d]", len(domains)), content, 300, 0)
			hasError := false
			for index, domain := range domains {
				progress := fmt.Sprintf("[%d/%d]", index+1, len(domains))
				w, err := whois.Lookup(domain.Name)
				content.Objects[0].(*widget.Label).SetText(fmt.Sprintf("Fetching WHOIS info of %s. %s", domain.Name, progress))
				window.SetTitle(fmt.Sprintf("Refreshing... %s", progress))
				if err != nil {
					content.Add(widget.NewLabel(fmt.Sprintf("%s Whois Lookup Error when refreshing domain %s: %v", progress, domain.Name, err)))
					hasError = true
				} else {
					if v := database.DB.Create(w); v.Error != nil {
						content.Add(widget.NewLabel(fmt.Sprintf("%s Database Error when refreshing domain %s: %v", progress, domain.Name, v.Error)))
						hasError = true
					}
				}
				whoises[index] = w
			}
			if hasError {
				content.Add(widget.NewButton("OK", func() {
					window.Close()
				}))
			} else {
				window.Close()
			}
			table.Refresh()
		}),
	)
}

func DomainList(inst *instance.Instance) fyne.CanvasObject {
	table := domainListTable(inst)
	toolbar := domainListToolbar(inst, table)
	return container.NewBorder(toolbar, nil, nil, nil, table)
}
