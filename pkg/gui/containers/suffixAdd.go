package containers

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	"DomainMan/pkg/models"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

func SuffixAdd(inst *instance.Instance, thisWindow fyne.Window, list *widget.List) fyne.CanvasObject {
	suffixName := widget.NewEntry()
	mode := widget.NewSelect([]string{"Whois Service", "Web"}, nil)
	mode.SetSelected("Whois Service")
	whoisServer := widget.NewEntry()
	var form *widget.Form
	handleError := func() {
		form.SubmitText = "Submit"
		form.Refresh()
		form.Enable()
	}
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Suffix", Widget: suffixName},
			{Text: "Mode", Widget: mode},
			{Text: "Whois Server", Widget: whoisServer},
		},
		OnSubmit: func() {
			form.SubmitText = "Adding..."
			form.Refresh()
			form.Disable()
			logrus.Infof("Adding suffix %s", suffixName.Text)
			suffix := models.Suffix{
				Name:        suffixName.Text,
				WhoisServer: whoisServer.Text,
			}
			switch mode.Selected {
			case "Whois Service":
				suffix.Mode = models.ModeWhois
			case "Web":
				suffix.Mode = models.ModeWeb
			default:
				suffix.Mode = models.ModeWhois
			}
			if v := database.DB.Create(&suffix); v.Error != nil {
				logrus.Errorf("Error adding suffix %s: %s", suffixName.Text, v.Error)
				utils.ShowSimpleWindow(inst, "Error", widget.NewLabel(v.Error.Error()))
				handleError()
				return
			}
			if list != nil {
				err := loadSuffixes()
				if err != nil {
					utils.ShowSimpleWindow(inst, "Database Error", widget.NewLabel(fmt.Sprintf("Error when loading suffixes: %v", err)))
				}
				list.Refresh()
			}
			thisWindow.Close()
		},
	}
	return form
}
