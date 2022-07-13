package containers

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	"DomainMan/pkg/models"
	"DomainMan/pkg/whois"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

func DomainAdd(inst *instance.Instance, thisWindow fyne.Window, list *widget.List) fyne.CanvasObject {
	domainName := widget.NewEntry()
	queryWhois := widget.NewCheck("", nil)
	queryWhois.SetChecked(true)
	var form *widget.Form
	handleError := func() {
		form.SubmitText = "Submit"
		form.Refresh()
		form.Enable()
	}
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Domain Name", Widget: domainName},
			{Text: "Query WHOIS after created", Widget: queryWhois},
		},
		OnSubmit: func() {
			form.SubmitText = "Adding..."
			form.Refresh()
			form.Disable()
			logrus.Infof("Adding domain %s", domainName.Text)
			domain := models.Domain{
				Name: domainName.Text,
			}
			if v := database.DB.Create(&domain); v.Error != nil {
				logrus.Errorf("Error adding domain %s: %s", domainName.Text, v.Error)
				utils.ShowSimpleWindow(inst, "Error", widget.NewLabel(v.Error.Error()))
				handleError()
				return
			}
			if queryWhois.Checked {
				form.SubmitText = "Querying WHOIS..."
				form.Refresh()
				logrus.Infof("Querying WHOIS for domain %s", domainName.Text)
				w, err := whois.Lookup(domain.Name)
				if err != nil {
					logrus.Errorf("Error querying WHOIS for domain %s: %s", domainName.Text, err)
					utils.ShowSimpleWindow(inst, "Error", widget.NewLabel(err.Error()))
					goto end
				}
				if v := database.DB.Create(w); v.Error != nil {
					logrus.Errorf("Error saving WHOIS for domain %s: %s", domainName.Text, v.Error)
					utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(v.Error.Error()))
					handleError()
					return
				}
			}
		end:
			if list != nil {
				err := loadDomains()
				if err != nil {
					utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(fmt.Sprintf("Error when loading domains: %v", err)))
				}
				list.Refresh()
			}
			thisWindow.Close()
		},
	}
	return form
}
