package containers

import (
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	"DomainMan/pkg/models"
	"DomainMan/pkg/whois"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/net/idna"
	"regexp"
)

func WhoisQuery(inst *instance.Instance) fyne.CanvasObject {
	domainNameEntry := widget.NewEntry()
	domainNameEntry.SetPlaceHolder("example.com")
	checkDomain := func(domain string) (string, error) {
		realDomain, err := idna.ToASCII(domain)
		if err != nil {
			return "", err
		}
		if !regexp.MustCompile("^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$").MatchString(realDomain) {
			return "", fmt.Errorf("invalid domain name")
		}
		return realDomain, nil
	}
	domainNameEntry.Validator = func(s string) error {
		_, err := checkDomain(s)
		return err
	}
	var queryButton *widget.Button
	doQuery := func() {
		domainNameEntry.Disable()
		queryButton.Disable()
		queryButton.SetText("Querying...")
		go func() {
			var w *models.Whois
			realDomainName, err := checkDomain(domainNameEntry.Text)
			defer func() {
				domainNameEntry.Enable()
				queryButton.Enable()
				queryButton.SetText("Query")
			}()
			if err != nil {
				utils.ShowSimpleWindow(inst, "Error", widget.NewLabel(err.Error()))
				return
			}
			w, err = whois.Lookup(realDomainName)
			if err != nil {
				utils.ShowSimpleWindow(inst, "Whois Lookup Error", widget.NewLabel(err.Error()))
			} else {
				utils.ShowSizedWindow(inst, fmt.Sprintf("Whois Info of %s", w.DomainName), WhoisDetail(inst, w), 600, 400)
			}
		}()
	}
	domainNameEntry.OnSubmitted = func(s string) {
		doQuery()
	}
	queryButton = widget.NewButton("Query", doQuery)
	title := canvas.NewText("WHOIS Query", theme.ForegroundColor())
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 24
	return container.NewPadded(container.NewVBox(
		title,
		domainNameEntry,
		queryButton,
	))
}
